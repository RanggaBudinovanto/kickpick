// Package jdsports scrapes jdsports.id, an authorized multi-brand sneaker
// retailer in Indonesia (part of JD Sports Fashion plc, sells genuine Nike/
// Adidas/etc. stock). This exists specifically because Nike's own site
// disallows crawling its product pages via robots.txt, and Adidas's own site
// is behind Akamai bot protection that rejects even basic requests — see
// PENDING.md. jdsports.id carries the same authentic products through its
// own storefront, with no such restriction, so this is the legitimate way to
// get Nike/Adidas price data without touching either brand's own site.
//
// robots.txt for jdsports.id (checked 2026-07-25) only disallows /search/,
// /checkout/, /profile/, /pickup-store/, /amp/ — product pages and the
// per-brand sitemap used here are unrestricted.
//
// The site is Next.js. Product listing pages are client-rendered (empty of
// product data in the raw HTML), but individual product pages are/aren't
// server-rendered — they embed the full product record, including price and
// stock, in a `__NEXT_DATA__` JSON script tag. Product URLs are enumerated
// from jdsports.id's own per-brand XML sitemap rather than the listing page,
// since that sitemap is a plain static file with no rendering involved.
package jdsports

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/kickpick/backend/internal/scraper"
)

const userAgent = "KickPickBot/1.0 (+https://kickpick.id/bot; price comparison crawler)"

// Adapter scrapes one brand's catalog on jdsports.id. brandSlug must match an
// existing row in the brands table (Section: Pipeline.Run resolves brand
// once per adapter and applies it to every scraped product), and
// sitemapURL must point at that brand's per-brand product sitemap, e.g.
// https://jdsports.id/sitemap/jdsport/nike-232.xml.
type Adapter struct {
	brandSlug  string
	sitemapURL string
	http       *http.Client
}

func New(brandSlug, sitemapURL string) *Adapter {
	return &Adapter{
		brandSlug:  brandSlug,
		sitemapURL: sitemapURL,
		http:       &http.Client{Timeout: 15 * time.Second},
	}
}

func (a *Adapter) BrandSlug() string { return a.brandSlug }

type sitemapURLSet struct {
	URLs []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

func (a *Adapter) Scrape(ctx context.Context) ([]scraper.ScrapedProduct, error) {
	if err := scraper.CheckRobotsAllowed(a.sitemapURL, userAgent); err != nil {
		return nil, fmt.Errorf("robots.txt check failed, refusing to scrape: %w", err)
	}

	productURLs, err := a.fetchProductURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("jdsports fetch sitemap: %w", err)
	}

	c := colly.NewCollector(colly.UserAgent(userAgent))
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*jdsports.id*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})
	c.SetRequestTimeout(20 * time.Second)

	var products []scraper.ScrapedProduct

	c.OnResponse(func(r *colly.Response) {
		p, err := parseProductPage(r.Body, a.brandSlug)
		if err != nil {
			return // skip pages that don't parse rather than fail the whole run
		}
		products = append(products, p)
	})

	for _, u := range productURLs {
		// Best-effort per URL: one bad product page shouldn't abort the run,
		// same rationale as the pipeline's own per-product error handling.
		_ = c.Visit(u)
	}
	c.Wait()

	return products, nil
}

func (a *Adapter) fetchProductURLs(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.sitemapURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := a.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching sitemap", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var set sitemapURLSet
	if err := xml.Unmarshal(body, &set); err != nil {
		return nil, fmt.Errorf("parse sitemap xml: %w", err)
	}

	urls := make([]string, 0, len(set.URLs))
	for _, u := range set.URLs {
		if u.Loc != "" {
			urls = append(urls, u.Loc)
		}
	}
	return urls, nil
}

var nextDataRe = regexp.MustCompile(`<script id="__NEXT_DATA__"[^>]*>(.*?)</script>`)
var ogImageRe = regexp.MustCompile(`<meta property="og:image" content="([^"]*)"`)

type nextDataProduct struct {
	Props struct {
		PageProps struct {
			DataProduct struct {
				SKU          string  `json:"sku"`
				Name         string  `json:"name"`
				Brand        string  `json:"brand"`
				Price        float64 `json:"price"`
				SpecialPrice float64 `json:"special_price"`
				StockStatus  int     `json:"stock_status"`
				URLKey       string  `json:"url_key"`
			} `json:"dataProduct"`
		} `json:"pageProps"`
	} `json:"props"`
}

func parseProductPage(body []byte, brandSlug string) (scraper.ScrapedProduct, error) {
	m := nextDataRe.FindSubmatch(body)
	if m == nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("__NEXT_DATA__ not found")
	}

	var data nextDataProduct
	if err := json.Unmarshal(m[1], &data); err != nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("parse __NEXT_DATA__: %w", err)
	}

	dp := data.Props.PageProps.DataProduct
	if dp.Name == "" || dp.URLKey == "" {
		return scraper.ScrapedProduct{}, fmt.Errorf("product data missing (likely not a product page)")
	}

	price := dp.SpecialPrice
	if price == 0 {
		price = dp.Price
	}
	if price == 0 {
		return scraper.ScrapedProduct{}, fmt.Errorf("no price found for %q", dp.Name)
	}

	imageURL := ""
	if im := ogImageRe.FindSubmatch(body); im != nil {
		imageURL = string(im[1])
	}

	productURL := "https://jdsports.id/product/" + dp.URLKey

	return scraper.ScrapedProduct{
		BrandSlug: brandSlug,
		Name:      dp.Name,
		Slug:      "jdsports-" + dp.URLKey,
		Category:  "sneakers",
		ImageURL:  imageURL,
		IsLimited: false,
		Offers: []scraper.ScrapedOffer{
			{
				StoreName:    "JD Sports Indonesia",
				StoreType:    "reseller",
				Price:        price,
				Currency:     "IDR",
				InStock:      dp.StockStatus == 1,
				AffiliateURL: productURL,
			},
		},
	}, nil
}
