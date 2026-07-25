// Package queue wraps Asynq (Redis-backed job queue) for the scrape and
// exchange-rate jobs, replacing the old in-process robfig/cron scheduler
// (PENDING.md). Because tasks are enqueued into Redis rather than fired by a
// local timer, running cmd/worker as more than one instance no longer
// duplicates work: whichever instance is free dequeues a given task, instead
// of every instance's own timer firing independently.
package queue

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/exchangerate"
	"github.com/kickpick/backend/internal/scraper"
	"github.com/kickpick/backend/internal/scraper/registry"
)

const (
	TypeScrapeRun         = "scrape:run"
	TypeExchangeRateFetch = "exchangerate:fetch"
)

// RedisConnOpt parses a redis://host:port URL (as used by REDIS_URL
// elsewhere in this codebase) into the connection option both the Asynq
// client and server need.
func RedisConnOpt(redisURL string) (asynq.RedisConnOpt, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, err
	}
	return opt, nil
}

func NewScrapeRunTask() *asynq.Task {
	return asynq.NewTask(TypeScrapeRun, nil)
}

func NewExchangeRateFetchTask() *asynq.Task {
	return asynq.NewTask(TypeExchangeRateFetch, nil)
}

// Handlers bundles the dependencies task handlers need so cmd/worker can
// register them on an asynq.ServeMux without reaching for globals.
type Handlers struct {
	pool     *pgxpool.Pool
	pipeline *scraper.Pipeline
}

func NewHandlers(pool *pgxpool.Pool, cfg *config.Config) *Handlers {
	return &Handlers{
		pool: pool,
		pipeline: scraper.NewPipeline(pool, scraper.PipelineConfig{
			AppURL:       cfg.AppURL,
			ResendAPIKey: cfg.ResendAPIKey,
			EmailFrom:    cfg.EmailFrom,
		}),
	}
}

// HandleScrapeRun mirrors the old scheduler's runAll: run every registered
// adapter, logging (not failing) per-adapter errors so one broken adapter
// doesn't block the others or trigger Asynq's automatic task retry — the
// next scheduled run is the natural retry, same as the old cron behavior.
func (h *Handlers) HandleScrapeRun(ctx context.Context, _ *asynq.Task) error {
	log.Println("[queue] starting scrape run")
	for _, adapter := range registry.All() {
		if err := h.pipeline.Run(ctx, adapter); err != nil {
			log.Printf("[queue] scrape run failed for %s: %v", adapter.BrandSlug(), err)
		}
	}
	log.Println("[queue] scrape run complete")
	return nil
}

func (h *Handlers) HandleExchangeRateFetch(ctx context.Context, _ *asynq.Task) error {
	log.Println("[queue] fetching exchange rate")
	if err := exchangerate.FetchAndStore(ctx, h.pool); err != nil {
		log.Printf("[queue] exchange rate fetch failed: %v", err)
		return nil
	}
	log.Println("[queue] exchange rate updated")
	return nil
}
