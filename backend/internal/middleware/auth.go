package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/kickpick/backend/internal/auth"
)

const LocalsUserID = "user_id"

func RequireAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "autentikasi diperlukan"})
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ParseAccessToken(secret, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "sesi berakhir, silakan login kembali"})
		}

		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "sesi tidak valid"})
		}

		c.Locals(LocalsUserID, uid)
		return c.Next()
	}
}

func UserIDFromContext(c *fiber.Ctx) (uuid.UUID, bool) {
	uid, ok := c.Locals(LocalsUserID).(uuid.UUID)
	return uid, ok
}
