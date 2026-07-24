package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/kickpick/backend/internal/auth"
	"github.com/kickpick/backend/internal/config"
)

const refreshCookieName = "refresh_token"

func setRefreshCookie(c *fiber.Ctx, rawToken string, cfg *config.Config) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    rawToken,
		Expires:  time.Now().Add(auth.RefreshTokenTTL),
		HTTPOnly: true,
		Secure:   cfg.AppEnv == "production",
		SameSite: "Strict",
		Path:     "/",
	})
}

func clearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "Strict",
		Path:     "/",
	})
}
