// Package testutil provides shared helpers for integration tests that need a
// real Postgres connection. Tests using this skip themselves (via t.Skip) when
// TEST_DATABASE_URL isn't set, so `go test ./...` still passes in environments
// without a database (e.g. a CI runner that hasn't provisioned Postgres yet).
package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RequireTestDB returns a connection pool to TEST_DATABASE_URL, or skips the
// calling test if that env var isn't set. Point it at a disposable database —
// tests that use this truncate tables they touch, but nothing here creates
// isolated schemas per test.
func RequireTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	return pool
}
