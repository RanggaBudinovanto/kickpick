// Command worker runs the scraping scheduler continuously (daily price
// updates, Section 10/19 PRD). Kept as a separate process from cmd/api so the
// HTTP server isn't blocked by or coupled to scrape run duration.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/scheduler"
)

// Default: once a day at 03:00 server time (low-traffic window).
const defaultCronSpec = "0 3 * * *"

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diset, cek file .env")
	}

	spec := os.Getenv("SCRAPE_CRON_SPEC")
	if spec == "" {
		spec = defaultCronSpec
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer pool.Close()

	cronJob, err := scheduler.Start(ctx, pool, spec)
	if err != nil {
		log.Fatalf("gagal start scheduler: %v", err)
	}

	log.Printf("worker berjalan, jadwal scrape: %q", spec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("worker shutting down...")
	cronCtx := cronJob.Stop()
	<-cronCtx.Done()
}
