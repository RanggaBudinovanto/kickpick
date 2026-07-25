// Package registry lists every brand adapter that's actually wired up. It's a
// separate package (rather than living in internal/scraper itself) because
// internal/scraper/compass imports internal/scraper for shared types, and
// registry needs to import both — putting this here avoids an import cycle.
//
// Section 24 PRD lists 6 candidate brands. Nike and Adidas's own sites
// are blocked (robots.txt disallow / Akamai bot protection — PENDING.md),
// so their data comes via authorized retailers (jdsports.id, planetsports.asia)
// instead of the brands' own sites, which is unblocked and legitimate.
package registry

import (
	"github.com/kickpick/backend/internal/config"
	"github.com/kickpick/backend/internal/scraper"
	"github.com/kickpick/backend/internal/scraper/compass"
	"github.com/kickpick/backend/internal/scraper/jdsports"
	"github.com/kickpick/backend/internal/scraper/plugo"
	"github.com/kickpick/backend/internal/scraper/shopee"
)

// All returns every adapter that's ready to actually run. The shopee adapter
// only gets included once its credentials are configured (PENDING.md) — it
// can't do anything yet, so registering it unconditionally would just log
// the same "not configured" error on every scheduled run for no benefit.
//
// Adidas is sourced from jdsports.id rather than planetsports.asia (which
// also carries it): planetsports.asia's robots.txt is malformed (Disallow
// rules with no preceding User-agent line — a bug on their end, not a
// deliberate block), and CheckRobotsAllowed fails closed on anything it
// can't parse. internal/scraper/planetsports still exists, ready to use if
// they ever fix their robots.txt.
func All(cfg *config.Config) []scraper.Adapter {
	adapters := []scraper.Adapter{
		compass.New(),
		jdsports.New("nike", "https://jdsports.id/sitemap/jdsport/nike-232.xml"),
		jdsports.New("adidas", "https://jdsports.id/sitemap/jdsport/adidas-223.xml"),
		jdsports.New("jordan", "https://jdsports.id/sitemap/jdsport/jordan-228.xml"),
		jdsports.New("new-balance", "https://jdsports.id/sitemap/jdsport/new-balance-230.xml"),
		jdsports.New("puma", "https://jdsports.id/sitemap/jdsport/puma-234.xml"),
		jdsports.New("asics", "https://jdsports.id/sitemap/jdsport/asics-987.xml"),
		jdsports.New("on", "https://jdsports.id/sitemap/jdsport/on-2166.xml"),
		jdsports.New("under-armour", "https://jdsports.id/sitemap/jdsport/under-armour-3008.xml"),
		jdsports.New("crocs", "https://jdsports.id/sitemap/jdsport/crocs-2062.xml"),
		jdsports.New("mizuno", "https://jdsports.id/sitemap/jdsport/mizuno-3028.xml"),
		plugo.New("geoffmax", "https://www.geoff-max.com", "Geoff Max Official Store"),
		plugo.New("brodo", "https://www.bro.do", "Brodo Official Store"),
	}

	shopeeCfg := shopee.Config{PartnerID: cfg.ShopeePartnerID, PartnerKey: cfg.ShopeePartnerKey}
	if shopeeCfg.Configured() {
		adapters = append(adapters, shopee.New(shopeeCfg))
	}

	return adapters
}
