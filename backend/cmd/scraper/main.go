// Command scraper runs all registered brand adapters once and persists the
// results. Used both for manual runs and as the job the scheduler (cmd/worker)
// triggers daily (Section 10 PRD: job scheduler for harga harian).
package main

import (
	"context"
	"log"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/scraper"
	"github.com/kickpick/backend/internal/scraper/registry"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diset, cek file .env")
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer pool.Close()

	pipeline := scraper.NewPipeline(pool)

	for _, adapter := range registry.All() {
		if err := pipeline.Run(ctx, adapter); err != nil {
			log.Printf("scrape run failed for %s: %v", adapter.BrandSlug(), err)
		}
	}

	log.Println("scrape run selesai")
}
