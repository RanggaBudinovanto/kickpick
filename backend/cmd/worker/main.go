// Command worker runs the scrape and exchange-rate jobs on a schedule
// (Section 10/19 PRD). Kept as a separate process from cmd/api so the HTTP
// server isn't blocked by or coupled to scrape run duration.
//
// Backed by Asynq (Redis-backed job queue) rather than an in-process cron
// timer, so a task enqueued at each cron tick is processed by whichever
// worker instance is free — see the RUN_SCHEDULER note below for how that
// interacts with running more than one instance of this binary.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"

	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/queue"
)

// Defaults: scrape once a day at 03:00, exchange rate once a day at 02:00
// (before the scrape, so a fresh rate is in place before the day's traffic).
const (
	defaultScrapeCronSpec = "0 3 * * *"
	exchangeRateCronSpec  = "0 2 * * *"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diset, cek file .env")
	}
	if cfg.RedisURL == "" {
		log.Fatal("REDIS_URL belum diset, cek file .env (worker sekarang butuh Redis untuk job queue)")
	}

	spec := os.Getenv("SCRAPE_CRON_SPEC")
	if spec == "" {
		spec = defaultScrapeCronSpec
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer pool.Close()

	redisOpt, err := queue.RedisConnOpt(cfg.RedisURL)
	if err != nil {
		log.Fatalf("REDIS_URL tidak valid: %v", err)
	}

	handlers := queue.NewHandlers(pool, cfg)
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeScrapeRun, handlers.HandleScrapeRun)
	mux.HandleFunc(queue.TypeExchangeRateFetch, handlers.HandleExchangeRateFetch)

	srv := asynq.NewServer(redisOpt, asynq.Config{Concurrency: 2})
	if err := srv.Start(mux); err != nil {
		log.Fatalf("gagal start asynq server: %v", err)
	}
	defer srv.Shutdown()

	client := asynq.NewClient(redisOpt)
	defer client.Close()

	// RUN_SCHEDULER lets a multi-replica deployment run exactly one
	// instance's Scheduler (the component that enqueues on cron ticks)
	// while every replica still runs a Server (the component that processes
	// queued tasks, above) for extra processing capacity. Without this, N
	// schedulers would each independently enqueue a duplicate task at every
	// tick — Asynq dedupes *processing* across instances via the shared
	// Redis queue, but does not dedupe *scheduling* on its own. Defaults to
	// true, which is correct for local dev and for the common case of
	// running a single worker replica.
	runScheduler := os.Getenv("RUN_SCHEDULER") != "false"
	var sched *asynq.Scheduler
	if runScheduler {
		sched = asynq.NewScheduler(redisOpt, nil)
		if _, err := sched.Register(spec, queue.NewScrapeRunTask()); err != nil {
			log.Fatalf("gagal daftar jadwal scrape: %v", err)
		}
		if _, err := sched.Register(exchangeRateCronSpec, queue.NewExchangeRateFetchTask()); err != nil {
			log.Fatalf("gagal daftar jadwal exchange rate: %v", err)
		}
		if err := sched.Start(); err != nil {
			log.Fatalf("gagal start asynq scheduler: %v", err)
		}

		// Enqueue both once immediately so a fresh deploy doesn't wait for
		// the first cron tick.
		if _, err := client.Enqueue(queue.NewScrapeRunTask()); err != nil {
			log.Printf("gagal enqueue scrape run awal: %v", err)
		}
		if _, err := client.Enqueue(queue.NewExchangeRateFetchTask()); err != nil {
			log.Printf("gagal enqueue exchange rate awal: %v", err)
		}
	}

	log.Printf("worker berjalan (scheduler=%v), jadwal scrape: %q, jadwal kurs: %q", runScheduler, spec, exchangeRateCronSpec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("worker shutting down...")
	if sched != nil {
		sched.Shutdown()
	}
}
