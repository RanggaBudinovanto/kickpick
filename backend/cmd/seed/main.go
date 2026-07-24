package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kickpick/backend/internal/auth"
	"github.com/kickpick/backend/internal/config"
	dbutil "github.com/kickpick/backend/internal/db"
	"github.com/kickpick/backend/internal/db/sqlc"
)

type seedBrand struct {
	name    string
	slug    string
	isLocal bool
}

type seedProduct struct {
	brandSlug string
	name      string
	slug      string
	category  string
	isLimited bool
	basePrice float64
}

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diset, cek file .env")
	}

	ctx := context.Background()
	pool, err := dbutil.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer pool.Close()

	q := sqlc.New(pool)

	brandIDs := seedBrands(ctx, pool)
	storeIDs := seedStores(ctx, pool)
	seedProducts(ctx, pool, brandIDs, storeIDs)
	seedDemoUser(ctx, q)
	seedExchangeRate(ctx, pool)

	log.Println("seed selesai")
}

func seedBrands(ctx context.Context, pool *pgxpool.Pool) map[string]uuid.UUID {
	brands := []seedBrand{
		{"Compass", "compass", true},
		{"Ventela", "ventela", true},
		{"Aerostreet", "aerostreet", true},
		{"Nike", "nike", false},
		{"Adidas", "adidas", false},
		{"Vans", "vans", false},
	}

	ids := make(map[string]uuid.UUID, len(brands))
	for _, b := range brands {
		id := uuid.New()
		ids[b.slug] = id
		mustExec(ctx, pool, `
			INSERT INTO brands (id, name, slug, is_local)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (slug) DO NOTHING
		`, id, b.name, b.slug, b.isLocal)
	}
	return ids
}

func seedStores(ctx context.Context, pool *pgxpool.Pool) map[string]uuid.UUID {
	stores := []struct {
		name    string
		typ     string
		network string
	}{
		{"Compass Official Store", "official", ""},
		{"Ventela Official Store", "official", ""},
		{"Shopee - Aerostreet Official", "marketplace", "shopee"},
		{"Tokopedia - Nike Official", "marketplace", "tokopedia"},
		{"Sport Station", "reseller", "involve_asia"},
		{"Planet Sports", "reseller", "involve_asia"},
	}

	ids := make(map[string]uuid.UUID, len(stores))
	for _, s := range stores {
		id := uuid.New()
		ids[s.name] = id
		mustExecPool(ctx, pool, `
			INSERT INTO stores (id, name, type, affiliate_network)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING
		`, id, s.name, s.typ, nullIfEmpty(s.network))
	}
	return ids
}

func seedProducts(ctx context.Context, pool *pgxpool.Pool, brandIDs map[string]uuid.UUID, storeIDs map[string]uuid.UUID) {
	products := []seedProduct{
		{"compass", "Compass Original Low", "compass-original-low", "lifestyle", false, 349000},
		{"compass", "Compass Court Legacy", "compass-court-legacy", "lifestyle", false, 399000},
		{"ventela", "Ventela Zoro Classic", "ventela-zoro-classic", "lifestyle", false, 289000},
		{"aerostreet", "Aerostreet Kirana", "aerostreet-kirana", "running", false, 259000},
		{"nike", "Nike Air Force 1 07", "nike-air-force-1-07", "lifestyle", false, 1799000},
		{"nike", "Nike Dunk Low Retro", "nike-dunk-low-retro", "lifestyle", true, 2199000},
		{"adidas", "Adidas Samba OG", "adidas-samba-og", "lifestyle", false, 1899000},
		{"adidas", "Adidas Ultraboost Light", "adidas-ultraboost-light", "running", false, 3299000},
		{"vans", "Vans Old Skool", "vans-old-skool", "lifestyle", false, 899000},
	}

	storeList := make([]uuid.UUID, 0, len(storeIDs))
	for _, id := range storeIDs {
		storeList = append(storeList, id)
	}

	rng := rand.New(rand.NewSource(42))

	for _, p := range products {
		brandID, ok := brandIDs[p.brandSlug]
		if !ok {
			continue
		}
		productID := uuid.New()
		mustExec(ctx, pool, `
			INSERT INTO products (id, brand_id, name, slug, category, is_limited)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (slug) DO NOTHING
		`, productID, brandID, p.name, p.slug, p.category, p.isLimited)

		mustExec(ctx, pool, `
			INSERT INTO product_images (id, product_id, url, sort_order)
			VALUES ($1, $2, $3, 0)
		`, uuid.New(), productID, "https://placehold.co/800x800.png?text="+p.slug)

		numOffers := 2 + rng.Intn(2)
		for i := 0; i < numOffers && i < len(storeList); i++ {
			storeID := storeList[(rng.Int())%len(storeList)]
			priceVariance := 1.0 + (rng.Float64()-0.5)*0.1
			price := p.basePrice * priceVariance
			mustExec(ctx, pool, `
				INSERT INTO product_offers (id, product_id, store_id, price, currency, in_stock, affiliate_url)
				VALUES ($1, $2, $3, $4, 'IDR', $5, $6)
			`, uuid.New(), productID, storeID, price, rng.Intn(10) > 1, "https://kickpick.id/go/"+p.slug)
		}

		for d := 30; d >= 0; d-- {
			date := time.Now().AddDate(0, 0, -d)
			drift := 1.0 + (rng.Float64()-0.5)*0.08
			mustExec(ctx, pool, `
				INSERT INTO price_history (id, product_id, store_id, price, recorded_date)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (product_id, store_id, recorded_date) DO NOTHING
			`, uuid.New(), productID, storeList[0], p.basePrice*drift, date)
		}
	}
}

func seedDemoUser(ctx context.Context, q *sqlc.Queries) {
	hash, err := auth.HashPassword("Password123")
	if err != nil {
		log.Fatalf("gagal hash password demo user: %v", err)
	}

	_, err = q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        "demo@kickpick.id",
		PasswordHash: hash,
		Name:         "Demo User",
	})
	if err != nil {
		log.Printf("demo user mungkin sudah ada, skip: %v", err)
	}
}

func seedExchangeRate(ctx context.Context, pool *pgxpool.Pool) {
	mustExecPool(ctx, pool, `
		INSERT INTO exchange_rates (id, base_currency, target_currency, rate, recorded_date)
		VALUES ($1, 'IDR', 'USD', 0.000062, CURRENT_DATE)
		ON CONFLICT (base_currency, target_currency, recorded_date) DO NOTHING
	`, uuid.New())
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func mustExec(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) {
	mustExecPool(ctx, pool, sql, args...)
}

func mustExecPool(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) {
	if _, err := pool.Exec(ctx, sql, args...); err != nil {
		log.Fatalf("seed query failed: %v\nsql: %s", err, sql)
	}
}
