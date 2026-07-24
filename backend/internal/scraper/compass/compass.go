// Package compass scrapes sepatucompass.com, the official Shopify (Hydrogen/Remix)
// storefront for the Compass sneaker brand.
//
// robots.txt for sepatucompass.com (checked 2026-07-24) only disallows
// /shop/dummy-archieve-sale and otherwise allows crawling ("Allow: /"), so
// fetching the public shop page is compliant with Section 14 PRD.
//
// The site server-renders its full product catalog into a `window.__remixContext`
// JS variable (Shopify Hydrogen's Remix loader data) rather than exposing a
// documented JSON API. The same product nodes appear repeatedly across multiple
// nested collection groupings in that blob, so instead of modeling Remix's exact
// (and undocumented, could change) schema, this adapter recursively walks the
// decoded JSON for any object shaped like a Shopify product node and dedupes by
// handle. This is more resilient to markup/schema churn than a brittle fixed path.
package compass

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/kickpick/backend/internal/scraper"
)

const (
	// /shop turned out to be an unstable route on this Hydrogen deployment (it
	// intermittently 404s). /collections/all-products is the same catalog and
	// has been reliable across repeated checks, so that's the one we scrape.
	shopURL   = "https://sepatucompass.com/collections/all-products"
	userAgent = "KickPickBot/1.0 (+https://kickpick.id/bot; price comparison crawler)"
)

type Adapter struct{}

func New() *Adapter {
	return &Adapter{}
}

func (a *Adapter) BrandSlug() string {
	return "compass"
}

func (a *Adapter) Scrape(ctx context.Context) ([]scraper.ScrapedProduct, error) {
	// Section 14 PRD: brands under scraping must be reviewed periodically, not
	// just checked once by hand during development — this makes that an
	// enforced pre-flight check on every run, so the adapter self-disables the
	// moment the site's robots.txt changes to disallow us, rather than someone
	// having to notice and update comments/code manually.
	if err := scraper.CheckRobotsAllowed(shopURL, userAgent); err != nil {
		return nil, fmt.Errorf("robots.txt check failed, refusing to scrape: %w", err)
	}

	c := colly.NewCollector(
		colly.UserAgent(userAgent),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*sepatucompass.com*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})
	c.SetRequestTimeout(20 * time.Second)

	var products []scraper.ScrapedProduct
	var scrapeErr error

	c.OnResponse(func(r *colly.Response) {
		nodes, err := extractProductNodes(r.Body)
		if err != nil {
			scrapeErr = err
			return
		}
		products = toScrapedProducts(nodes)
	})

	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("compass scrape request failed: %w", err)
	})

	if err := c.Visit(shopURL); err != nil {
		return nil, fmt.Errorf("compass scrape visit failed: %w", err)
	}
	c.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}
	return products, nil
}

const remixContextMarker = "window.__remixContext = "

// productNode mirrors the fields we care about from a Shopify Storefront API
// product node, as embedded in the page's SSR'd loader data.
type productNode struct {
	Title            string `json:"title"`
	Handle           string `json:"handle"`
	AvailableForSale bool   `json:"availableForSale"`
	PriceRange       struct {
		MinVariantPrice struct {
			Amount       string `json:"amount"`
			CurrencyCode string `json:"currencyCode"`
		} `json:"minVariantPrice"`
	} `json:"priceRange"`
	Images struct {
		Nodes []struct {
			URL string `json:"url"`
		} `json:"nodes"`
	} `json:"images"`
}

func extractProductNodes(body []byte) (map[string]productNode, error) {
	content := string(body)

	start := strings.Index(content, remixContextMarker)
	if start == -1 {
		return nil, fmt.Errorf("remixContext marker not found in response (page structure may have changed)")
	}
	jsonStart := start + len(remixContextMarker)

	jsonEnd, err := findJSONObjectEnd(content, jsonStart)
	if err != nil {
		return nil, fmt.Errorf("failed to locate end of remixContext JSON: %w", err)
	}

	var raw interface{}
	if err := json.Unmarshal([]byte(content[jsonStart:jsonEnd]), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse remixContext JSON: %w", err)
	}

	found := make(map[string]productNode)
	walkForProducts(raw, found)
	return found, nil
}

// findJSONObjectEnd returns the index just past the closing brace of the JSON
// object starting at content[start], by tracking brace depth and string state
// (including escapes). The trailing script content after the assignment isn't
// stable across pages/deploys, so this is more robust than string-searching
// for a specific end marker.
func findJSONObjectEnd(content string, start int) (int, error) {
	if start >= len(content) || content[start] != '{' {
		return 0, fmt.Errorf("expected '{' at position %d", start)
	}

	depth := 0
	inString := false
	escaped := false

	for i := start; i < len(content); i++ {
		c := content[i]

		if inString {
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}

		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1, nil
			}
		}
	}

	return 0, fmt.Errorf("unterminated JSON object")
}

// walkForProducts recursively scans decoded JSON for objects shaped like a
// Shopify product node (has title + handle + priceRange), deduping by handle.
func walkForProducts(node interface{}, found map[string]productNode) {
	switch v := node.(type) {
	case map[string]interface{}:
		if isProductNode(v) {
			var p productNode
			if b, err := json.Marshal(v); err == nil {
				if err := json.Unmarshal(b, &p); err == nil && p.Handle != "" {
					found[p.Handle] = p
				}
			}
		}
		for _, val := range v {
			walkForProducts(val, found)
		}
	case []interface{}:
		for _, item := range v {
			walkForProducts(item, found)
		}
	}
}

func isProductNode(v map[string]interface{}) bool {
	_, hasTitle := v["title"]
	_, hasHandle := v["handle"]
	_, hasPriceRange := v["priceRange"]
	return hasTitle && hasHandle && hasPriceRange
}

func toScrapedProducts(nodes map[string]productNode) []scraper.ScrapedProduct {
	products := make([]scraper.ScrapedProduct, 0, len(nodes))
	for handle, n := range nodes {
		price, err := strconv.ParseFloat(n.PriceRange.MinVariantPrice.Amount, 64)
		if err != nil {
			continue
		}

		imageURL := ""
		if len(n.Images.Nodes) > 0 {
			imageURL = n.Images.Nodes[0].URL
		}

		currency := n.PriceRange.MinVariantPrice.CurrencyCode
		if currency == "" {
			currency = "IDR"
		}

		products = append(products, scraper.ScrapedProduct{
			BrandSlug: "compass",
			Name:      n.Title,
			Slug:      "compass-" + handle,
			Category:  "lifestyle",
			ImageURL:  imageURL,
			IsLimited: false,
			Offers: []scraper.ScrapedOffer{
				{
					StoreName:    "Compass Official Store",
					StoreType:    "official",
					Price:        price,
					Currency:     currency,
					InStock:      n.AvailableForSale,
					AffiliateURL: "https://sepatucompass.com/shop/" + handle,
				},
			},
		})
	}
	return products
}
