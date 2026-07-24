package scheduler

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/exchangerate"
	"github.com/kickpick/backend/internal/scraper"
	"github.com/kickpick/backend/internal/scraper/registry"
)

// Start runs the scrape pipeline for every registered adapter on the given
// cron schedule (in addition to once immediately on startup, so a fresh
// deploy doesn't wait a full day for its first data). Section 10/19 PRD: job
// scheduler for daily price updates.
//
// Still robfig/cron in-process rather than Redis-backed Asynq (PENDING.md) —
// Redis is available now (used for rate limiting), but switching the job
// queue itself is a separate, larger change not done yet. Fine for a single
// backend instance; two instances would each fire their own cron independently.
func Start(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config, spec string) (*cron.Cron, error) {
	pipeline := scraper.NewPipeline(pool, scraper.PipelineConfig{
		AppURL:       cfg.AppURL,
		ResendAPIKey: cfg.ResendAPIKey,
		EmailFrom:    cfg.EmailFrom,
	})
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

	// Runs before the scrape (02:00 vs 03:00) so a fresh rate is in place
	// before /api/exchange-rate and the frontend's currency toggle are used
	// during the day. A failed fetch just leaves the previous day's rate in
	// place (GetLatestExchangeRate) rather than blocking anything.
	runExchangeRate := func() {
		log.Println("[scheduler] fetching daily exchange rate")
		if err := exchangerate.FetchAndStore(ctx, pool); err != nil {
			log.Printf("[scheduler] exchange rate fetch failed: %v", err)
			return
		}
		log.Println("[scheduler] exchange rate updated")
	}

	c := cron.New()
	if _, err := c.AddFunc(spec, runAll); err != nil {
		return nil, err
	}
	if _, err := c.AddFunc("0 2 * * *", runExchangeRate); err != nil {
		return nil, err
	}

	go runAll()
	go runExchangeRate()

	c.Start()
	return c, nil
}
