package scraper

import "context"

// ScrapedOffer is one store's listing for a scraped product. Most direct-site
// scrapers (like Compass) produce exactly one offer per product; adapters that
// aggregate multiple marketplaces for the same product would produce more.
type ScrapedOffer struct {
	StoreName    string
	StoreType    string // "official", "marketplace", "reseller"
	Price        float64
	Currency     string
	InStock      bool
	AffiliateURL string
}

// ScrapedProduct is the normalized shape every brand adapter must produce,
// regardless of how the source site structures its own data.
type ScrapedProduct struct {
	BrandSlug string
	Name      string
	Slug      string // stable per-brand identifier (e.g. Shopify handle)
	Category  string
	ImageURL  string
	IsLimited bool
	Offers    []ScrapedOffer
}

// Adapter is the interface every per-brand scraper implements. Section 10 PRD:
// modular scraping per brand, Colly for static sites / chromedp for JS-heavy ones.
// The interface doesn't care which one a given adapter uses internally.
type Adapter interface {
	BrandSlug() string
	Scrape(ctx context.Context) ([]ScrapedProduct, error)
}
