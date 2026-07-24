package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	token, err := GenerateAccessToken(secret, userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := ParseAccessToken(secret, token)
	if err != nil {
		t.Fatalf("ParseAccessToken failed: %v", err)
	}
	if claims.UserID != userID.String() {
		t.Errorf("claims.UserID = %q, want %q", claims.UserID, userID.String())
	}
}

func TestParseAccessTokenRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateAccessToken("secret-a", userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	if _, err := ParseAccessToken("secret-b", token); err == nil {
		t.Error("ParseAccessToken should reject a token signed with a different secret")
	}
}

func TestParseAccessTokenRejectsGarbage(t *testing.T) {
	if _, err := ParseAccessToken("secret", "not-a-real-token"); err == nil {
		t.Error("ParseAccessToken should reject malformed tokens")
	}
}

func TestGenerateOpaqueTokenIsUniqueAndHashConsistent(t *testing.T) {
	raw1, hash1, err := GenerateOpaqueToken()
	if err != nil {
		t.Fatalf("GenerateOpaqueToken failed: %v", err)
	}
	raw2, hash2, err := GenerateOpaqueToken()
	if err != nil {
		t.Fatalf("GenerateOpaqueToken failed: %v", err)
	}

	if raw1 == raw2 {
		t.Error("two calls to GenerateOpaqueToken produced the same raw token")
	}
	if hash1 == hash2 {
		t.Error("two calls to GenerateOpaqueToken produced the same hash")
	}
	if HashToken(raw1) != hash1 {
		t.Error("HashToken(raw) should reproduce the same hash returned by GenerateOpaqueToken")
	}
}
