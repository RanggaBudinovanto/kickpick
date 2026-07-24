// Package exchangerate fetches the daily IDR->USD rate from a free,
// no-API-key-required source (Frankfurter, backed by ECB reference rates)
// and stores it in the exchange_rates table so the frontend's currency
// toggle (Section 10/19 PRD) reflects a real, current rate instead of the
// static seed value.
package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

const frankfurterURL = "https://api.frankfurter.app/latest?from=IDR&to=USD"

type frankfurterResponse struct {
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// FetchAndStore fetches today's IDR->USD rate and upserts it into
// exchange_rates, keyed by (base_currency, target_currency, recorded_date).
func FetchAndStore(ctx context.Context, pool *pgxpool.Pool) error {
	rate, err := fetchIDRToUSD(ctx)
	if err != nil {
		return fmt.Errorf("fetch exchange rate: %w", err)
	}

	queries := sqlc.New(pool)
	if err := queries.UpsertExchangeRate(ctx, sqlc.UpsertExchangeRateParams{
		BaseCurrency:   "IDR",
		TargetCurrency: "USD",
		Rate:           db.Float64ToNumeric(rate),
		RecordedDate:   db.Date(time.Now()),
	}); err != nil {
		return fmt.Errorf("store exchange rate: %w", err)
	}
	return nil
}

func fetchIDRToUSD(ctx context.Context) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, frankfurterURL, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status %d from frankfurter.app", resp.StatusCode)
	}

	var body frankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}

	rate, ok := body.Rates["USD"]
	if !ok {
		return 0, fmt.Errorf("USD rate missing from response")
	}
	return rate, nil
}
