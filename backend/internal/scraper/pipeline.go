package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
	"github.com/kickpick/backend/internal/email"
)

type Pipeline struct {
	Queries *sqlc.Queries
	AppURL  string
	Email   *email.Client
	http    *http.Client
}

type PipelineConfig struct {
	// AppURL is the frontend's base URL (e.g. https://kickpick.id), used to
	// call its /api/revalidate route after a price update so ISR-cached
	// product pages (Section 19 PRD) don't sit on stale prices for the full
	// cache window, and to build links in restock/price-drop emails. Pass ""
	// to skip revalidation (e.g. in tests).
	AppURL string
	// ResendAPIKey/EmailFrom configure restock/price-drop alert emails
	// (Section 17 PRD E04/E05). Pass an empty ResendAPIKey to log instead of
	// send (see internal/email.Client).
	ResendAPIKey string
	EmailFrom    string
}

func NewPipeline(pool *pgxpool.Pool, cfg PipelineConfig) *Pipeline {
	return &Pipeline{
		Queries: sqlc.New(pool),
		AppURL:  cfg.AppURL,
		Email:   email.NewClient(cfg.ResendAPIKey, cfg.EmailFrom),
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

// Run scrapes with the given adapter and persists results: upserts the product
// record, replaces its primary image, upserts each offer (creating stores as
// needed), and appends today's price to price_history. Section 10/19 PRD: this
// is what the daily scheduled job calls per brand.
func (p *Pipeline) Run(ctx context.Context, adapter Adapter) error {
	brand, err := p.Queries.GetBrandBySlug(ctx, adapter.BrandSlug())
	if err != nil {
		return fmt.Errorf("brand %q not found in database (seed it first): %w", adapter.BrandSlug(), err)
	}

	products, err := adapter.Scrape(ctx)
	if err != nil {
		return fmt.Errorf("scrape failed for brand %q: %w", adapter.BrandSlug(), err)
	}

	log.Printf("[scraper] %s: scraped %d products", adapter.BrandSlug(), len(products))

	var errCount int
	for _, sp := range products {
		if err := p.persistProduct(ctx, brand.ID, sp); err != nil {
			log.Printf("[scraper] %s: failed to persist %q: %v", adapter.BrandSlug(), sp.Name, err)
			errCount++
		}
	}

	log.Printf("[scraper] %s: persisted %d/%d products (%d errors)", adapter.BrandSlug(), len(products)-errCount, len(products), errCount)
	return nil
}

func (p *Pipeline) persistProduct(ctx context.Context, brandID pgtype.UUID, sp ScrapedProduct) error {
	product, err := p.Queries.UpsertProduct(ctx, sqlc.UpsertProductParams{
		BrandID:   brandID,
		Name:      sp.Name,
		Slug:      sp.Slug,
		Category:  sp.Category,
		IsLimited: sp.IsLimited,
	})
	if err != nil {
		return fmt.Errorf("upsert product: %w", err)
	}

	if err := p.Queries.ReplaceProductImage(ctx, sqlc.ReplaceProductImageParams{
		ProductID: product.ID,
		Url:       sp.ImageURL,
	}); err != nil {
		return fmt.Errorf("replace image: %w", err)
	}

	for _, offer := range sp.Offers {
		if err := p.persistOffer(ctx, product, offer); err != nil {
			return fmt.Errorf("persist offer from %q: %w", offer.StoreName, err)
		}
	}

	p.revalidateProduct(sp.Slug)

	return nil
}

func (p *Pipeline) revalidateProduct(slug string) {
	if p.AppURL == "" {
		return
	}

	body, err := json.Marshal(map[string]string{"tag": "product:" + slug})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, p.AppURL+"/api/revalidate", bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Best-effort: a revalidation failure shouldn't fail the scrape run, the
	// ISR window (1 hour) is still a safe fallback if this doesn't land.
	if res, err := p.http.Do(req); err == nil {
		_ = res.Body.Close()
	}
}

func (p *Pipeline) persistOffer(ctx context.Context, product sqlc.Product, offer ScrapedOffer) error {
	store, err := p.getOrCreateStore(ctx, offer.StoreName, offer.StoreType)
	if err != nil {
		return fmt.Errorf("get/create store: %w", err)
	}

	priceNumeric := dbutil.Float64ToNumeric(offer.Price)

	existing, err := p.Queries.GetOfferByProductAndStore(ctx, sqlc.GetOfferByProductAndStoreParams{
		ProductID: product.ID,
		StoreID:   store.ID,
	})
	hadPreviousOffer := true
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("check existing offer: %w", err)
		}
		hadPreviousOffer = false
		if _, err := p.Queries.CreateProductOffer(ctx, sqlc.CreateProductOfferParams{
			ProductID:    product.ID,
			StoreID:      store.ID,
			Price:        priceNumeric,
			Currency:     offer.Currency,
			InStock:      offer.InStock,
			AffiliateUrl: offer.AffiliateURL,
		}); err != nil {
			return fmt.Errorf("create offer: %w", err)
		}
	} else {
		if err := p.Queries.UpdateProductOffer(ctx, sqlc.UpdateProductOfferParams{
			ProductID:    product.ID,
			StoreID:      store.ID,
			Price:        priceNumeric,
			Currency:     offer.Currency,
			InStock:      offer.InStock,
			AffiliateUrl: offer.AffiliateURL,
		}); err != nil {
			return fmt.Errorf("update offer: %w", err)
		}
	}

	if err := p.Queries.UpsertPriceHistoryToday(ctx, sqlc.UpsertPriceHistoryTodayParams{
		ProductID: product.ID,
		StoreID:   store.ID,
		Price:     priceNumeric,
	}); err != nil {
		return fmt.Errorf("upsert price history: %w", err)
	}

	// Alert subscribers (Section 17 PRD E04/E05) only on a genuine state
	// change against a PREVIOUSLY KNOWN offer — a brand-new offer (first time
	// we've ever seen this store carry this product) is not a "restock" or
	// "price drop" from the subscriber's point of view.
	if hadPreviousOffer {
		wasOutOfStock := !existing.InStock
		justRestocked := wasOutOfStock && offer.InStock
		oldPrice := dbutil.NumericToFloat64(existing.Price)
		priceDropped := offer.InStock && offer.Price < oldPrice

		if justRestocked {
			p.notifySubscribers(ctx, product, notificationRestock)
		} else if priceDropped {
			p.notifySubscribers(ctx, product, notificationPriceDrop)
		}
	}

	return nil
}

type notificationKind int

const (
	notificationRestock notificationKind = iota
	notificationPriceDrop
)

// notifySubscribers creates an in-app notification and sends an email
// (Section 17 PRD templates E04/E05) to every user with an active wishlist
// alert on this product. Best-effort per recipient — one failure doesn't
// stop the others or fail the scrape run.
func (p *Pipeline) notifySubscribers(ctx context.Context, product sqlc.Product, kind notificationKind) {
	subscribers, err := p.Queries.ListActiveAlertSubscribersForProduct(ctx, product.ID)
	if err != nil {
		log.Printf("[scraper] failed to list alert subscribers for %q: %v", product.Name, err)
		return
	}
	if len(subscribers) == 0 {
		return
	}

	productURL := fmt.Sprintf("%s/id/produk/%s", p.AppURL, product.Slug)

	var title, body, emailSubject, emailHTML string
	switch kind {
	case notificationRestock:
		title = "Restock!"
		body = fmt.Sprintf("%s baru saja tersedia lagi.", product.Name)
		emailSubject, emailHTML = email.RestockAlertEmail(product.Name, productURL)
	case notificationPriceDrop:
		title = "Harga turun"
		body = fmt.Sprintf("Harga %s baru saja turun.", product.Name)
		emailSubject, emailHTML = email.PriceDropAlertEmail(product.Name, productURL)
	}

	notifType := map[notificationKind]string{
		notificationRestock:   "restock",
		notificationPriceDrop: "price_drop",
	}[kind]

	for _, sub := range subscribers {
		if _, err := p.Queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
			UserID:    sub.UserID,
			Type:      notifType,
			Title:     title,
			Body:      body,
			ActionUrl: dbutil.Text(productURL),
		}); err != nil {
			log.Printf("[scraper] failed to create notification for user %s: %v", sub.Email, err)
			continue
		}

		if err := p.Email.Send(sub.Email, emailSubject, emailHTML); err != nil {
			log.Printf("[scraper] failed to send alert email to %s: %v", sub.Email, err)
		}
	}

	log.Printf("[scraper] notified %d subscriber(s) for %q (%s)", len(subscribers), product.Name, notifType)
}

func (p *Pipeline) getOrCreateStore(ctx context.Context, name, storeType string) (sqlc.Store, error) {
	store, err := p.Queries.GetStoreByName(ctx, name)
	if err == nil {
		return store, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return sqlc.Store{}, err
	}

	return p.Queries.CreateStore(ctx, sqlc.CreateStoreParams{
		Name:             name,
		Type:             storeType,
		AffiliateNetwork: dbutil.Text(""),
	})
}
