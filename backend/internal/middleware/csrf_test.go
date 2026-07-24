package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRequireCSRFHeader(t *testing.T) {
	app := fiber.New()
	app.Post("/action", RequireCSRFHeader(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	t.Run("rejects request without header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/action", nil)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed: %v", err)
		}
		if res.StatusCode != fiber.StatusForbidden {
			t.Errorf("status = %d, want 403", res.StatusCode)
		}
	})

	t.Run("rejects request with wrong header value", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/action", nil)
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed: %v", err)
		}
		if res.StatusCode != fiber.StatusForbidden {
			t.Errorf("status = %d, want 403", res.StatusCode)
		}
	})

	t.Run("accepts request with correct header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/action", nil)
		req.Header.Set("X-Requested-With", "kickpick")
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test failed: %v", err)
		}
		if res.StatusCode != fiber.StatusOK {
			t.Errorf("status = %d, want 200", res.StatusCode)
		}
	})
}
