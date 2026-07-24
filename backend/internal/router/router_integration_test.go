package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/router"
	"github.com/kickpick/backend/internal/testutil"
)

// These exercise the actual HTTP stack (routing, CSRF middleware, auth
// middleware, real Postgres) end to end for the critical paths called out in
// Section 20 PRD: register -> login -> logout, and an authenticated action
// blocked without a valid session.
func TestAuthFlow_RegisterLoginProfile(t *testing.T) {
	pool := testutil.RequireTestDB(t)

	cfg := &config.Config{
		AppEnv:            "test",
		JWTAccessSecret:   "test-access-secret",
		JWTRefreshSecret:  "test-refresh-secret",
		CORSAllowedOrigin: "http://localhost:3000",
		EmailFrom:         "KickPick <noreply@kickpick.id>",
	}
	app := router.New(cfg, pool)

	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())

	// Register
	registerBody, _ := json.Marshal(map[string]string{
		"email":            email,
		"password":         "Password123",
		"confirm_password": "Password123",
		"name":             "Test User",
	})
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	if res.StatusCode != 201 {
		t.Fatalf("register status = %d, want 201", res.StatusCode)
	}

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "Password123",
	})
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("login status = %d, want 200", res.StatusCode)
	}

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&loginResp); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}
	if loginResp.AccessToken == "" {
		t.Fatal("login response did not include an access_token")
	}

	// Profile without token: must be rejected.
	req = httptest.NewRequest("GET", "/api/profile", nil)
	res, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("unauthenticated profile request failed: %v", err)
	}
	if res.StatusCode != 401 {
		t.Errorf("unauthenticated profile status = %d, want 401", res.StatusCode)
	}

	// Profile with token: must succeed and return the registered email.
	req = httptest.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	res, err = app.Test(req, -1)
	if err != nil {
		t.Fatalf("authenticated profile request failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("authenticated profile status = %d, want 200", res.StatusCode)
	}

	var profileResp struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(res.Body).Decode(&profileResp); err != nil {
		t.Fatalf("failed to decode profile response: %v", err)
	}
	if profileResp.Email != email {
		t.Errorf("profile email = %q, want %q", profileResp.Email, email)
	}
}

func TestAuthFlow_DuplicateEmailRejected(t *testing.T) {
	pool := testutil.RequireTestDB(t)

	cfg := &config.Config{
		JWTAccessSecret:   "test-access-secret",
		CORSAllowedOrigin: "http://localhost:3000",
		EmailFrom:         "KickPick <noreply@kickpick.id>",
	}
	app := router.New(cfg, pool)

	email := fmt.Sprintf("dup-%d@example.com", time.Now().UnixNano())
	body, _ := json.Marshal(map[string]string{
		"email":            email,
		"password":         "Password123",
		"confirm_password": "Password123",
		"name":             "Test User",
	})

	for i, wantStatus := range []int{201, 409} {
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("register request %d failed: %v", i, err)
		}
		if res.StatusCode != wantStatus {
			t.Errorf("attempt %d: status = %d, want %d", i, res.StatusCode, wantStatus)
		}
	}
}

func TestCSRFProtection_BlocksStateChangingRequestWithoutHeader(t *testing.T) {
	pool := testutil.RequireTestDB(t)

	cfg := &config.Config{
		JWTAccessSecret:   "test-access-secret",
		CORSAllowedOrigin: "http://localhost:3000",
	}
	app := router.New(cfg, pool)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("logout request failed: %v", err)
	}
	if res.StatusCode != 403 {
		t.Errorf("logout without CSRF header status = %d, want 403", res.StatusCode)
	}
}
