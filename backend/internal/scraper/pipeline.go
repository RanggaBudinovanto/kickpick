package scraper

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

type Pipeline struct {
	Queries *sqlc.Queries
}

func NewPipeline(pool *pgxpool.Pool) *Pipeline {
	return &Pipeline{Queries: sqlc.New(pool)}
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
		if err := p.persistOffer(ctx, product.ID, offer); err != nil {
			return fmt.Errorf("persist offer from %q: %w", offer.StoreName, err)
		}
	}

	return nil
}

func (p *Pipeline) persistOffer(ctx context.Context, productID pgtype.UUID, offer ScrapedOffer) error {
	store, err := p.getOrCreateStore(ctx, offer.StoreName, offer.StoreType)
	if err != nil {
		return fmt.Errorf("get/create store: %w", err)
	}

	priceNumeric := dbutil.Float64ToNumeric(offer.Price)

	existing, err := p.Queries.GetOfferByProductAndStore(ctx, sqlc.GetOfferByProductAndStoreParams{
		ProductID: productID,
		StoreID:   store.ID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("check existing offer: %w", err)
		}
		if _, err := p.Queries.CreateProductOffer(ctx, sqlc.CreateProductOfferParams{
			ProductID:    productID,
			StoreID:      store.ID,
			Price:        priceNumeric,
			Currency:     offer.Currency,
			InStock:      offer.InStock,
			AffiliateUrl: offer.AffiliateURL,
		}); err != nil {
			return fmt.Errorf("create offer: %w", err)
		}
	} else {
		_ = existing
		if err := p.Queries.UpdateProductOffer(ctx, sqlc.UpdateProductOfferParams{
			ProductID:    productID,
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
		ProductID: productID,
		StoreID:   store.ID,
		Price:     priceNumeric,
	}); err != nil {
		return fmt.Errorf("upsert price history: %w", err)
	}

	return nil
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
