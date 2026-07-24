package middleware

import "github.com/gofiber/fiber/v2"

// RequireCSRFHeader is a defense-in-depth CSRF check for state-changing requests.
// Refresh/logout rely on the SameSite=Strict cookie for CSRF protection already;
// this middleware additionally requires a custom header that a cross-site <form>
// submission cannot set without triggering a CORS preflight our strict origin
// policy would reject. Section 14 PRD.
func RequireCSRFHeader() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("X-Requested-With") != "kickpick" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permintaan ditolak"})
		}
		return c.Next()
	}
}
