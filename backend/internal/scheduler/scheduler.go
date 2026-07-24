package scheduler

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"

	"github.com/kickpick/backend/internal/scraper"
	"github.com/kickpick/backend/internal/scraper/registry"
)

// Start runs the scrape pipeline for every registered adapter on the given
// cron schedule (in addition to once immediately on startup, so a fresh
// deploy doesn't wait a full day for its first data). Section 10/19 PRD: job
// scheduler for daily price updates.
//
// No Redis/Asynq here (PENDING.md) — Asynq needs Redis, which isn't installed
// in this environment. robfig/cron running in-process is a reasonable
// single-instance interim; revisit if KickPick ever runs multiple backend
// instances, since two instances would each fire their own cron independently.
func Start(ctx context.Context, pool *pgxpool.Pool, spec string) (*cron.Cron, error) {
	pipeline := scraper.NewPipeline(pool)
	adapters := registry.All()

	runAll := func() {
		log.Println("[scheduler] starting scheduled scrape run")
		for _, adapter := range adapters {
			if err := pipeline.Run(ctx, adapter); err != nil {
				log.Printf("[scheduler] scrape run failed for %s: %v", adapter.BrandSlug(), err)
			}
		}
		log.Println("[scheduler] scheduled scrape run complete")
	}

	c := cron.New()
	if _, err := c.AddFunc(spec, runAll); err != nil {
		return nil, err
	}

	go runAll()

	c.Start()
	return c, nil
}
