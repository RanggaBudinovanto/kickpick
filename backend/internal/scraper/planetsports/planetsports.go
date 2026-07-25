// Package planetsports scrapes planetsports.asia, an authorized multi-brand
// sports retailer in Indonesia (part of MAP Group, sells genuine Adidas/
// Nike/etc. stock). Exists for the same reason as internal/scraper/jdsports:
// Adidas's own site is behind Akamai bot protection, so this authorized
// retailer's own storefront is the legitimate path to real Adidas price data
// without touching adidas.co.id at all. See PENDING.md.
//
// robots.txt for planetsports.asia (checked 2026-07-25) is a standard
// Magento default — disallows admin/checkout/customer-account/tag/review
// paths, not product pages.
//
// The site is classic server-rendered Magento: category pages list products
// directly in HTML (paginated via ?p=N), and product pages embed price,
// name, and stock status directly in the markup — no JS rendering needed,
// unlike jdsports.id's listing pages.
package planetsports

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/kickpick/backend/internal/scraper"
)

const userAgent = "KickPickBot/1.0 (+https://kickpick.id/bot; price comparison crawler)"

// maxCategoryPages caps pagination so a change in site behavior (e.g. an
// infinite redirect loop, or a much larger catalog than expected) can't turn
// into an unbounded scrape.
const maxCategoryPages = 30

// Adapter scrapes one brand's category on planetsports.asia. brandSlug must
// match an existing row in the brands table, and categoryURL is that
// brand's category listing page, e.g.
// https://www.planetsports.asia/adidas.html.
type Adapter struct {
	brandSlug   string
	categoryURL string
}

func New(brandSlug, categoryURL string) *Adapter {
	return &Adapter{brandSlug: brandSlug, categoryURL: categoryURL}
}

func (a *Adapter) BrandSlug() string { return a.brandSlug }

var productLinkRe = regexp.MustCompile(`href="(https://www\.planetsports\.asia/[a-z0-9-]+\.html)"`)

func (a *Adapter) Scrape(ctx context.Context) ([]scraper.ScrapedProduct, error) {
	if err := scraper.CheckRobotsAllowed(a.categoryURL, userAgent); err != nil {
		return nil, fmt.Errorf("robots.txt check failed, refusing to scrape: %w", err)
	}

	productURLs, err := a.fetchProductURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("planetsports fetch category pages: %w", err)
	}

	c := colly.NewCollector(colly.UserAgent(userAgent))
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*planetsports.asia*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})
	c.SetRequestTimeout(20 * time.Second)

	var products []scraper.ScrapedProduct

	c.OnResponse(func(r *colly.Response) {
		p, err := parseProductPage(r.Body, a.brandSlug)
		if err != nil {
			return
		}
		products = append(products, p)
	})

	for _, u := range productURLs {
		_ = c.Visit(u)
	}
	c.Wait()

	return products, nil
}

// fetchProductURLs pages through the category listing (?p=1, ?p=2, ...)
// collecting product links, stopping once a page yields nothing new (the
// site loops back to page 1 content past the last real page) or the safety
// cap is hit.
func (a *Adapter) fetchProductURLs(ctx context.Context) ([]string, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}
	seen := make(map[string]bool)
	var urls []string

	for page := 1; page <= maxCategoryPages; page++ {
		pageURL := fmt.Sprintf("%s?p=%d", a.categoryURL, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		body := make([]byte, 0)
		buf := make([]byte, 32*1024)
		for {
			n, rerr := resp.Body.Read(buf)
			if n > 0 {
				body = append(body, buf[:n]...)
			}
			if rerr != nil {
				break
			}
		}
		resp.Body.Close()

		matches := productLinkRe.FindAllSubmatch(body, -1)
		newOnThisPage := 0
		for _, m := range matches {
			u := string(m[1])
			if !seen[u] {
				seen[u] = true
				urls = append(urls, u)
				newOnThisPage++
			}
		}

		if newOnThisPage == 0 {
			break
		}

		time.Sleep(2 * time.Second) // same rate limit as the product-page crawl
	}

	return urls, nil
}

var (
	priceRe      = regexp.MustCompile(`<span class="price">Rp\.\s*([0-9.,]+)</span>`)
	ogTitleRe    = regexp.MustCompile(`<meta property="og:title" content="([^"]*)"`)
	ogImageRe    = regexp.MustCompile(`<meta property="og:image" content="([^"]*)"`)
	skuRe        = regexp.MustCompile(`"sku":"([^"]*)"`)
	stockDimRe   = regexp.MustCompile(`"dimension4":"([^"]*)"`)
	canonicalURL = regexp.MustCompile(`<link rel="canonical" href="([^"]*)"`)
)

func parseProductPage(body []byte, brandSlug string) (scraper.ScrapedProduct, error) {
	titleMatch := ogTitleRe.FindSubmatch(body)
	if titleMatch == nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("not a product page (no og:title)")
	}
	name := string(titleMatch[1])

	priceMatch := priceRe.FindSubmatch(body)
	if priceMatch == nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("no price found for %q", name)
	}
	price, err := parseIDRPrice(string(priceMatch[1]))
	if err != nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("parse price for %q: %w", name, err)
	}

	sku := ""
	if m := skuRe.FindSubmatch(body); m != nil {
		sku = string(m[1])
	}
	if sku == "" {
		return scraper.ScrapedProduct{}, fmt.Errorf("no sku found for %q", name)
	}

	imageURL := ""
	if m := ogImageRe.FindSubmatch(body); m != nil {
		imageURL = string(m[1])
	}

	productURL := ""
	if m := canonicalURL.FindSubmatch(body); m != nil {
		productURL = string(m[1])
	}

	inStock := true // fail open on stock status: absence of the dimension isn't proof it's out of stock
	if m := stockDimRe.FindSubmatch(body); m != nil {
		inStock = strings.EqualFold(string(m[1]), "in stock")
	}

	return scraper.ScrapedProduct{
		BrandSlug: brandSlug,
		Name:      name,
		Slug:      "planetsports-" + strings.ToLower(sku),
		Category:  "sneakers",
		ImageURL:  imageURL,
		IsLimited: false,
		Offers: []scraper.ScrapedOffer{
			{
				StoreName:    "Planet Sports",
				StoreType:    "reseller",
				Price:        price,
				Currency:     "IDR",
				InStock:      inStock,
				AffiliateURL: productURL,
			},
		},
	}, nil
}

// parseIDRPrice parses Indonesian-formatted price text like "3.500.000"
// (dot as thousands separator, no decimal) into a float64.
func parseIDRPrice(s string) (float64, error) {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(s, 64)
}
