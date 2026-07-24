package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/kickpick/backend/internal/auth"
)

func newTestApp(secret string) *fiber.App {
	app := fiber.New()
	app.Get("/protected", RequireAuth(secret), func(c *fiber.Ctx) error {
		uid, ok := UserIDFromContext(c)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"user_id": uid.String()})
	})
	return app
}

func TestRequireAuth_RejectsMissingHeader(t *testing.T) {
	app := newTestApp("secret")
	req := httptest.NewRequest("GET", "/protected", nil)

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}

func TestRequireAuth_RejectsMalformedHeader(t *testing.T) {
	app := newTestApp("secret")
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "NotBearer sometoken")

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}

func TestRequireAuth_RejectsTokenSignedWithDifferentSecret(t *testing.T) {
	app := newTestApp("secret-a")
	token, err := auth.GenerateAccessToken("secret-b", uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}

func TestRequireAuth_AcceptsValidToken(t *testing.T) {
	secret := "secret"
	userID := uuid.New()
	app := newTestApp(secret)

	token, err := auth.GenerateAccessToken(secret, userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Errorf("status = %d, want 200", res.StatusCode)
	}
}
