package exchangerate

import (
	"context"
	"testing"
	"time"

	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/testutil"
)

// TestFetchAndStore hits the real frankfurter.app API and a real test
// database (skipped if TEST_DATABASE_URL isn't set) — this is deliberately
// not mocked, since the thing worth verifying is that the live API's
// response shape still matches what fetchIDRToUSD expects and that the
// upsert actually lands a row GetLatestExchangeRate can read back.
func TestFetchAndStore(t *testing.T) {
	pool := testutil.RequireTestDB(t)
	ctx := context.Background()

	if err := FetchAndStore(ctx, pool); err != nil {
		t.Fatalf("FetchAndStore failed: %v", err)
	}

	queries := sqlc.New(pool)
	rate, err := queries.GetLatestExchangeRate(ctx, sqlc.GetLatestExchangeRateParams{
		BaseCurrency:   "IDR",
		TargetCurrency: "USD",
	})
	if err != nil {
		t.Fatalf("GetLatestExchangeRate failed after FetchAndStore: %v", err)
	}

	if !rate.RecordedDate.Valid || rate.RecordedDate.Time.Format("2006-01-02") != time.Now().Format("2006-01-02") {
		t.Errorf("expected recorded_date to be today, got %v", rate.RecordedDate)
	}

	f, err := rate.Rate.Float64Value()
	if err != nil || !f.Valid || f.Float64 <= 0 {
		t.Errorf("expected a positive stored rate, got %v (err: %v)", rate.Rate, err)
	}
}
