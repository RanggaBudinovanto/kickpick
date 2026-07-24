package handler

import (
	"context"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/auth"
	"github.com/kickpick/backend/internal/config"
	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/email"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type AuthHandler struct {
	Queries *sqlc.Queries
	Pool    *pgxpool.Pool
	Cfg     *config.Config
	Email   *email.Client
}

func NewAuthHandler(pool *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		Queries: sqlc.New(pool),
		Pool:    pool,
		Cfg:     cfg,
		Email:   email.NewClient(cfg.ResendAPIKey, cfg.EmailFrom),
	}
}

type registerRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Name            string `json:"name"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if !emailRegex.MatchString(req.Email) {
		return badRequest(c, "format email tidak valid")
	}
	if len(req.Password) < 8 {
		return badRequest(c, "password minimal 8 karakter")
	}
	if req.Password != req.ConfirmPassword {
		return badRequest(c, "konfirmasi password tidak cocok")
	}
	if req.Name == "" {
		return badRequest(c, "nama wajib diisi")
	}

	ctx := c.Context()

	if _, err := h.Queries.GetUserByEmail(ctx, req.Email); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email sudah digunakan"})
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return serverError(c, err)
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return serverError(c, err)
	}

	user, err := h.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
	})
	if err != nil {
		return serverError(c, err)
	}

	rawToken, tokenHash, err := auth.GenerateOpaqueToken()
	if err != nil {
		return serverError(c, err)
	}

	_, err = h.Queries.CreateEmailVerificationToken(ctx, sqlc.CreateEmailVerificationTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: dbutil.Timestamptz(time.Now().Add(24 * time.Hour)),
	})
	if err != nil {
		return serverError(c, err)
	}

	go func() {
		welcomeSubject, welcomeHTML := email.WelcomeEmail(req.Name)
		if err := h.Email.Send(req.Email, welcomeSubject, welcomeHTML); err != nil {
			log.Printf("failed to send welcome email to %s: %v", req.Email, err)
		}
		verifySubject, verifyHTML := email.VerifyEmail(h.Cfg.AppURL, rawToken)
		if err := h.Email.Send(req.Email, verifySubject, verifyHTML); err != nil {
			log.Printf("failed to send verification email to %s: %v", req.Email, err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "registrasi berhasil, cek email untuk verifikasi",
	})
}

type verifyEmailRequest struct {
	Token string `json:"token"`
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var req verifyEmailRequest
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return badRequest(c, "token tidak valid")
	}

	ctx := c.Context()
	tokenHash := auth.HashToken(req.Token)

	record, err := h.Queries.GetEmailVerificationToken(ctx, tokenHash)
	if err != nil {
		return badRequest(c, "token tidak valid atau kadaluarsa")
	}

	if err := h.Queries.SetUserEmailVerified(ctx, record.UserID); err != nil {
		return serverError(c, err)
	}
	if err := h.Queries.MarkEmailVerificationTokenUsed(ctx, record.ID); err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"message": "email berhasil diverifikasi"})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	ctx := c.Context()

	user, err := h.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email atau password salah"})
	}

	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email atau password salah"})
	}

	uid := dbutil.ToUUID(user.ID)
	logAudit(ctx, h.Queries, &uid, "login", c.IP())

	return h.issueSession(c, ctx, uid)
}

func (h *AuthHandler) issueSession(c *fiber.Ctx, ctx context.Context, uid uuid.UUID) error {
	accessToken, err := auth.GenerateAccessToken(h.Cfg.JWTAccessSecret, uid)
	if err != nil {
		return serverError(c, err)
	}

	rawRefresh, refreshHash, err := auth.GenerateOpaqueToken()
	if err != nil {
		return serverError(c, err)
	}

	_, err = h.Queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    dbutil.UUID(uid),
		TokenHash: refreshHash,
		ExpiresAt: dbutil.Timestamptz(time.Now().Add(auth.RefreshTokenTTL)),
	})
	if err != nil {
		return serverError(c, err)
	}

	setRefreshCookie(c, rawRefresh, h.Cfg)

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	rawRefresh := c.Cookies("refresh_token")
	if rawRefresh == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "sesi berakhir, silakan login kembali"})
	}

	ctx := c.Context()
	tokenHash := auth.HashToken(rawRefresh)

	record, err := h.Queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "sesi berakhir, silakan login kembali"})
	}

	accessToken, err := auth.GenerateAccessToken(h.Cfg.JWTAccessSecret, dbutil.ToUUID(record.UserID))
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"access_token": accessToken})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	rawRefresh := c.Cookies("refresh_token")
	if rawRefresh != "" {
		ctx := c.Context()
		tokenHash := auth.HashToken(rawRefresh)
		if record, err := h.Queries.GetRefreshToken(ctx, tokenHash); err == nil {
			uid := dbutil.ToUUID(record.UserID)
			logAudit(ctx, h.Queries, &uid, "logout", c.IP())
		}
		_ = h.Queries.RevokeRefreshToken(ctx, tokenHash)
	}
	clearRefreshCookie(c)
	return c.JSON(fiber.Map{"message": "logout berhasil"})
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req forgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	ctx := c.Context()
	user, err := h.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Tidak membocorkan apakah email terdaftar atau tidak.
		return c.JSON(fiber.Map{"message": "jika email terdaftar, instruksi reset sudah dikirim"})
	}

	rawToken, tokenHash, err := auth.GenerateOpaqueToken()
	if err != nil {
		return serverError(c, err)
	}

	_, err = h.Queries.CreatePasswordResetToken(ctx, sqlc.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: dbutil.Timestamptz(time.Now().Add(1 * time.Hour)),
	})
	if err != nil {
		return serverError(c, err)
	}

	go func() {
		subject, html := email.ResetPasswordEmail(h.Cfg.AppURL, rawToken)
		if err := h.Email.Send(req.Email, subject, html); err != nil {
			log.Printf("failed to send reset password email to %s: %v", req.Email, err)
		}
	}()

	return c.JSON(fiber.Map{"message": "jika email terdaftar, instruksi reset sudah dikirim"})
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req resetPasswordRequest
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return badRequest(c, "token tidak valid")
	}
	if len(req.NewPassword) < 8 {
		return badRequest(c, "password minimal 8 karakter")
	}

	ctx := c.Context()
	tokenHash := auth.HashToken(req.Token)

	record, err := h.Queries.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return badRequest(c, "token tidak valid atau kadaluarsa")
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return serverError(c, err)
	}

	if err := h.Queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           record.UserID,
		PasswordHash: hash,
	}); err != nil {
		return serverError(c, err)
	}

	if err := h.Queries.MarkPasswordResetTokenUsed(ctx, record.ID); err != nil {
		return serverError(c, err)
	}

	if err := h.Queries.RevokeAllUserRefreshTokens(ctx, record.UserID); err != nil {
		return serverError(c, err)
	}

	uid := dbutil.ToUUID(record.UserID)
	logAudit(ctx, h.Queries, &uid, "password_reset", c.IP())

	return c.JSON(fiber.Map{"message": "password berhasil direset, silakan login kembali"})
}

func badRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
}

func serverError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "terjadi kesalahan di server"})
}
