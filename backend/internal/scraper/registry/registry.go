// Package registry lists every brand adapter that's actually wired up. It's a
// separate package (rather than living in internal/scraper itself) because
// internal/scraper/compass imports internal/scraper for shared types, and
// registry needs to import both — putting this here avoids an import cycle.
//
// Section 24 PRD lists 6 candidate brands, but only Compass has a working
// adapter today — see PENDING.md for why the other 5 aren't scrapable yet
// (site down, no on-site catalog, or blocked by robots.txt/bot protection).
package registry

import (
	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/scraper"
	"github.com/kickpick/backend/internal/scraper/compass"
	"github.com/kickpick/backend/internal/scraper/shopee"
)

// All returns every adapter that's ready to actually run. The shopee adapter
// only gets included once its credentials are configured (PENDING.md) — it
// can't do anything yet, so registering it unconditionally would just log
// the same "not configured" error on every scheduled run for no benefit.
func All(cfg *config.Config) []scraper.Adapter {
	adapters := []scraper.Adapter{
		compass.New(),
	}

	shopeeCfg := shopee.Config{PartnerID: cfg.ShopeePartnerID, PartnerKey: cfg.ShopeePartnerKey}
	if shopeeCfg.Configured() {
		adapters = append(adapters, shopee.New(shopeeCfg))
	}

	return adapters
}
