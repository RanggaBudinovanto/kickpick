// Package shopee is a scaffold for a Shopee Affiliate adapter — the
// affiliate network chosen (2026-07-25) for brands that don't have their own
// scrapable storefront (Aerostreet, Vans, and others per the findings table
// in PENDING.md).
//
// Deliberately NOT wired to Shopee's real API yet. Shopee Affiliate
// Program/Open Platform requires registering a partner account, going
// through their approval process, and getting a Partner ID + Partner Key
// issued to that specific account — none of which exists yet. Guessing at
// request/response shapes without a verified, current copy of Shopee's API
// docs would risk silently shipping something that looks wired up but is
// actually wrong against their real contract. So this stops at: compiles,
// implements the same scraper.Adapter interface as compass, and returns a
// clear "not configured" error — ready to fill in once credentials and docs
// are in hand. See PENDING.md for what's needed to activate it.
package shopee

import (
	"context"
	"errors"

	"github.com/kickpick/backend/internal/scraper"
)

// Config holds the credentials issued after Shopee approves a partner
// account. Populate from env vars (SHOPEE_PARTNER_ID / SHOPEE_PARTNER_KEY)
// once they exist — there's nothing to fill in before that.
type Config struct {
	PartnerID  string
	PartnerKey string
}

// Configured reports whether enough credentials are present to attempt a
// real request. Callers (e.g. registry.All) should skip registering this
// adapter entirely while it's false, rather than run it daily just to log
// the same "not configured" error.
func (c Config) Configured() bool {
	return c.PartnerID != "" && c.PartnerKey != ""
}

type Adapter struct {
	cfg Config
}

func New(cfg Config) *Adapter {
	return &Adapter{cfg: cfg}
}

func (a *Adapter) BrandSlug() string {
	return "shopee"
}

var ErrNotConfigured = errors.New("shopee adapter: PartnerID/PartnerKey belum diset, lihat PENDING.md")

func (a *Adapter) Scrape(ctx context.Context) ([]scraper.ScrapedProduct, error) {
	if !a.cfg.Configured() {
		return nil, ErrNotConfigured
	}

	// TODO once credentials + a verified copy of Shopee's current Affiliate
	// Open API docs are available:
	//   - Request signing (Shopee's spec: HMAC-SHA256 over partner id +
	//     timestamp + request payload, exact format per their docs).
	//   - This is marketplace search, not a single brand's own catalog like
	//     compass — needs a decision on which keyword(s)/shop(s) to query
	//     per target brand (Aerostreet, Vans, ...), since Shopee doesn't
	//     expose a "give me everything from this brand" endpoint the same
	//     way a brand's own storefront does.
	//   - Cross-source dedup against any brand that later also gets its own
	//     direct-site adapter (PENDING.md notes this as a separate, not yet
	//     built, concern).
	return nil, errors.New("shopee adapter: belum diimplementasi — tunggu kredensial partner account dan dokumentasi API resmi Shopee")
}
