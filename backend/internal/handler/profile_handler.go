package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/auth"
	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/middleware"
)

type ProfileHandler struct {
	Queries *sqlc.Queries
}

func NewProfileHandler(pool *pgxpool.Pool) *ProfileHandler {
	return &ProfileHandler{Queries: sqlc.New(pool)}
}

func userToDTO(u sqlc.User) fiber.Map {
	return fiber.Map{
		"id":                 dbutil.ToUUID(u.ID).String(),
		"email":              u.Email,
		"name":               u.Name,
		"onboarding_focus":   u.OnboardingFocus.String,
		"preferred_language": u.PreferredLanguage,
		"preferred_currency": u.PreferredCurrency,
		"email_verified":     u.EmailVerified,
	}
}

func (h *ProfileHandler) Get(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	user, err := h.Queries.GetUserByID(c.Context(), dbutil.UUID(uid))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user tidak ditemukan"})
	}

	return c.JSON(userToDTO(user))
}

type updateProfileRequest struct {
	Name              string `json:"name"`
	OnboardingFocus   string `json:"onboarding_focus"`
	PreferredLanguage string `json:"preferred_language"`
	PreferredCurrency string `json:"preferred_currency"`
}

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return badRequest(c, "nama wajib diisi")
	}
	if req.PreferredLanguage != "id" && req.PreferredLanguage != "en" {
		req.PreferredLanguage = "id"
	}
	if req.PreferredCurrency != "IDR" && req.PreferredCurrency != "USD" {
		req.PreferredCurrency = "IDR"
	}

	user, err := h.Queries.UpdateUserProfile(c.Context(), sqlc.UpdateUserProfileParams{
		ID:                dbutil.UUID(uid),
		Name:              req.Name,
		OnboardingFocus:   dbutil.Text(req.OnboardingFocus),
		PreferredLanguage: req.PreferredLanguage,
		PreferredCurrency: req.PreferredCurrency,
	})
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(userToDTO(user))
}

type deleteAccountRequest struct {
	Password string `json:"password"`
}

func (h *ProfileHandler) Delete(c *fiber.Ctx) error {
	uid, _ := middleware.UserIDFromContext(c)

	var req deleteAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	user, err := h.Queries.GetUserByID(c.Context(), dbutil.UUID(uid))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user tidak ditemukan"})
	}

	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "password salah"})
	}

	if err := h.Queries.SoftDeleteUser(c.Context(), dbutil.UUID(uid)); err != nil {
		return serverError(c, err)
	}
	if err := h.Queries.RevokeAllUserRefreshTokens(c.Context(), dbutil.UUID(uid)); err != nil {
		return serverError(c, err)
	}

	logAudit(c.Context(), h.Queries, &uid, "account_deleted", c.IP())

	clearRefreshCookie(c)

	return c.JSON(fiber.Map{"message": "akun berhasil dihapus"})
}
