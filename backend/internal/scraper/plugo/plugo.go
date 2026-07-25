// Package plugo scrapes storefronts built on Plugo (an Indonesian e-commerce
// SaaS platform, Nuxt-based). Geoff Max (geoff-max.com) and Brodo (bro.do)
// both run on it with an identical template, so one adapter parametrized by
// domain covers both rather than duplicating near-identical code.
//
// robots.txt for both sites (checked 2026-07-25) is fully permissive
// (`Allow: /`). Product listing pages are client-rendered, but individual
// product pages embed price, name, and image directly in server-rendered
// HTML — product URLs are enumerated from the site's own sitemap.xml
// (a sitemapindex pointing at static/categories/products sub-sitemaps)
// rather than the listing page.
package plugo

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/kickpick/backend/internal/scraper"
)

const userAgent = "KickPickBot/1.0 (+https://kickpick.id/bot; price comparison crawler)"

// Adapter scrapes one Plugo-based storefront. baseURL is the site root (no
// trailing slash, e.g. "https://www.geoff-max.com"), and storeName is used
// as the offer's store label (e.g. "Geoff Max Official Store").
type Adapter struct {
	brandSlug string
	baseURL   string
	storeName string
	http      *http.Client
}

func New(brandSlug, baseURL, storeName string) *Adapter {
	return &Adapter{
		brandSlug: brandSlug,
		baseURL:   baseURL,
		storeName: storeName,
		http:      &http.Client{Timeout: 15 * time.Second},
	}
}

func (a *Adapter) BrandSlug() string { return a.brandSlug }

type sitemapIndex struct {
	Sitemaps []struct {
		Loc string `xml:"loc"`
	} `xml:"sitemap"`
}

type sitemapURLSet struct {
	URLs []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

func (a *Adapter) Scrape(ctx context.Context) ([]scraper.ScrapedProduct, error) {
	sitemapRoot := a.baseURL + "/sitemap.xml"
	if err := scraper.CheckRobotsAllowed(sitemapRoot, userAgent); err != nil {
		return nil, fmt.Errorf("robots.txt check failed, refusing to scrape: %w", err)
	}

	productURLs, err := a.fetchProductURLs(ctx, sitemapRoot)
	if err != nil {
		return nil, fmt.Errorf("plugo fetch sitemap: %w", err)
	}

	c := colly.NewCollector(colly.UserAgent(userAgent))
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*" + hostGlob(a.baseURL) + "*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})
	c.SetRequestTimeout(20 * time.Second)

	var products []scraper.ScrapedProduct

	c.OnResponse(func(r *colly.Response) {
		p, err := a.parseProductPage(r.Body)
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

func hostGlob(baseURL string) string {
	h := strings.TrimPrefix(baseURL, "https://")
	h = strings.TrimPrefix(h, "http://")
	h = strings.TrimPrefix(h, "www.")
	return h
}

func (a *Adapter) fetchXML(ctx context.Context, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := a.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, v)
}

// fetchProductURLs resolves the sitemapindex, then reads every sub-sitemap
// whose loc looks like a product listing (contains "/products") — the
// catalog can be split across multiple numbered files (products-0001.xml,
// products-0002.xml, ...), so this doesn't assume there's only one.
func (a *Adapter) fetchProductURLs(ctx context.Context, sitemapRoot string) ([]string, error) {
	var index sitemapIndex
	if err := a.fetchXML(ctx, sitemapRoot, &index); err != nil {
		return nil, err
	}

	var urls []string
	for _, sm := range index.Sitemaps {
		if !strings.Contains(sm.Loc, "/products") {
			continue
		}
		var set sitemapURLSet
		if err := a.fetchXML(ctx, sm.Loc, &set); err != nil {
			return nil, fmt.Errorf("fetch product sitemap %s: %w", sm.Loc, err)
		}
		for _, u := range set.URLs {
			if u.Loc != "" {
				urls = append(urls, u.Loc)
			}
		}
	}
	return urls, nil
}

var (
	ogTitleRe    = regexp.MustCompile(`<meta property="og:title" content="([^"]*)"`)
	ogImageRe    = regexp.MustCompile(`<meta property="og:image" content="([^"]*)"`)
	priceRe      = regexp.MustCompile(`text-h4-medium[^"]*"><div>Rp\s*([0-9,]+)</div>`)
	canonicalRe  = regexp.MustCompile(`<link rel="canonical" href="([^"]*)"`)
)

func (a *Adapter) parseProductPage(body []byte) (scraper.ScrapedProduct, error) {
	titleMatch := ogTitleRe.FindSubmatch(body)
	if titleMatch == nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("not a product page (no og:title)")
	}
	name := cleanTitle(string(titleMatch[1]))

	if isNonFootwear(name) {
		return scraper.ScrapedProduct{}, fmt.Errorf("skipping non-footwear product %q", name)
	}

	priceMatch := priceRe.FindSubmatch(body)
	if priceMatch == nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("no price found for %q", name)
	}
	price, err := strconv.ParseFloat(strings.ReplaceAll(string(priceMatch[1]), ",", ""), 64)
	if err != nil {
		return scraper.ScrapedProduct{}, fmt.Errorf("parse price for %q: %w", name, err)
	}

	imageURL := ""
	if m := ogImageRe.FindSubmatch(body); m != nil {
		imageURL = string(m[1])
	}

	productURL := ""
	if m := canonicalRe.FindSubmatch(body); m != nil {
		productURL = string(m[1])
	}

	slug := slugify(a.brandSlug + "-" + name)

	return scraper.ScrapedProduct{
		BrandSlug: a.brandSlug,
		Name:      name,
		Slug:      slug,
		Category:  "sneakers",
		ImageURL:  imageURL,
		IsLimited: false,
		Offers: []scraper.ScrapedOffer{
			{
				StoreName: a.storeName,
				StoreType: "official",
				Price:     price,
				Currency:  "IDR",
				// Stock state isn't reliably readable from static markup on
				// this template (sold-out styling is applied client-side),
				// so this fails open rather than guessing wrong in either
				// direction — see PENDING.md.
				InStock:      true,
				AffiliateURL: productURL,
			},
		},
	}, nil
}

// cleanTitle strips the site-name suffix Plugo's default template appends to
// every og:title, e.g. "Gavin Stance White Gum - GMX | GEOFF MAX Footwear -
// Official Webstore" -> "Gavin Stance White Gum".
func cleanTitle(title string) string {
	if i := strings.Index(title, " - "); i > 0 {
		return title[:i]
	}
	return title
}

// nonFootwearKeywords catches apparel/accessories these multi-category
// stores also sell (Brodo especially — wallets, bags, clothing) so KickPick,
// a shoe price-comparison site, doesn't end up with non-shoe listings.
var nonFootwearKeywords = []string{
	"crewneck", "hoodie", "jaket", "kaos", "kemeja", "celana",
	"dompet", "tas", "topi", "belt", "wallet", "bag", "shirt",
	"socks", "kaus kaki", "sarung tangan", "gloves", "sabuk",
}

func isNonFootwear(name string) bool {
	lower := strings.ToLower(name)
	for _, kw := range nonFootwearKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = nonSlugChars.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
