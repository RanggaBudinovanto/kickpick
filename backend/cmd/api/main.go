package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/router"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool := mustConnect(ctx, cfg)
	defer func() {
		if pool != nil {
			pool.Close()
		}
	}()

	app := router.New(cfg, pool)

	log.Printf("KickPick API listening on :%s (env=%s)", cfg.AppPort, cfg.AppEnv)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func mustConnect(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	if cfg.DatabaseURL == "" {
		log.Println("DATABASE_URL not set, starting without database connection")
		return nil
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return pool
}
