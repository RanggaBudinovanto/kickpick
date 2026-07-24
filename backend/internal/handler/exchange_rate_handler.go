package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

type ExchangeRateHandler struct {
	Queries *sqlc.Queries
}

func NewExchangeRateHandler(pool *pgxpool.Pool) *ExchangeRateHandler {
	return &ExchangeRateHandler{Queries: sqlc.New(pool)}
}

// GetIDRToUSD returns the most recently recorded IDR->USD rate (Section 10/19
// PRD: harga dasar disimpan IDR, USD adalah hasil konversi tampilan).
func (h *ExchangeRateHandler) GetIDRToUSD(c *fiber.Ctx) error {
	rate, err := h.Queries.GetLatestExchangeRate(c.Context(), sqlc.GetLatestExchangeRateParams{
		BaseCurrency:   "IDR",
		TargetCurrency: "USD",
	})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "kurs belum tersedia"})
	}

	return c.JSON(fiber.Map{
		"base_currency":   rate.BaseCurrency,
		"target_currency": rate.TargetCurrency,
		"rate":            dbutil.NumericToFloat64(rate.Rate),
		"recorded_date":   rate.RecordedDate.Time.Format("2006-01-02"),
	})
}
