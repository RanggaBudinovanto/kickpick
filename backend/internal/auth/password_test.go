package auth

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("Password123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "Password123" {
		t.Fatal("password was not hashed")
	}
	if !VerifyPassword(hash, "Password123") {
		t.Error("VerifyPassword should succeed with the correct password")
	}
	if VerifyPassword(hash, "wrong-password") {
		t.Error("VerifyPassword should fail with an incorrect password")
	}
}

func TestHashPasswordProducesDifferentHashesForSameInput(t *testing.T) {
	h1, _ := HashPassword("Password123")
	h2, _ := HashPassword("Password123")
	if h1 == h2 {
		t.Error("bcrypt should salt each hash uniquely, got identical hashes")
	}
}
