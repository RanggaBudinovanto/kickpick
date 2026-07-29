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
		{"Jordan", "jordan", false},
		{"New Balance", "new-balance", false},
		{"Puma", "puma", false},
		{"Asics", "asics", false},
		{"On", "on", false},
		{"Under Armour", "under-armour", false},
		{"Crocs", "crocs", false},
		{"Mizuno", "mizuno", false},
		{"Geoff Max", "geoffmax", true},
		{"Brodo", "brodo", true},
	}

	ids := make(map[string]uuid.UUID, len(brands))
	for _, b := range brands {
		logoURL := "/brands/" + b.slug + ".svg"
		// Use RETURNING id so we always get the real DB id, even on conflict.
		var id uuid.UUID
		row := pool.QueryRow(ctx, `
			INSERT INTO brands (id, name, slug, is_local, logo_url)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (slug) DO UPDATE SET logo_url = EXCLUDED.logo_url
			RETURNING id
		`, uuid.New(), b.name, b.slug, b.isLocal, logoURL)
		if err := row.Scan(&id); err != nil {
			log.Fatalf("seedBrands failed for %s: %v", b.slug, err)
		}
		ids[b.slug] = id
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
	products := []struct {
		seedProduct
		imageURL string
	}{
		{seedProduct{"compass", "Compass Original Low", "compass-original-low", "lifestyle", false, 349000}, "https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"compass", "Compass Court Legacy", "compass-court-legacy", "lifestyle", false, 399000}, "https://images.unsplash.com/photo-1560769629-975ec94e6a86?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"ventela", "Ventela Zoro Classic", "ventela-zoro-classic", "lifestyle", false, 289000}, "https://images.unsplash.com/photo-1549298916-b41d501d3772?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"aerostreet", "Aerostreet Kirana", "aerostreet-kirana", "running", false, 259000}, "https://images.unsplash.com/photo-1542291026-7eec264c27ff?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"nike", "Nike Air Force 1 07", "nike-air-force-1-07", "lifestyle", false, 1799000}, "https://images.unsplash.com/photo-1595950653106-6c9ebd614d3a?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"nike", "Nike Dunk Low Retro", "nike-dunk-low-retro", "lifestyle", true, 2199000}, "https://images.unsplash.com/photo-1600185365926-3a2ce3cdb9eb?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"nike", "Nike React Infinity Run", "nike-react-infinity-run", "running", false, 2499000}, "https://images.unsplash.com/photo-1584735935682-2f2b69dff9d2?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"adidas", "Adidas Samba OG", "adidas-samba-og", "lifestyle", false, 1899000}, "https://images.unsplash.com/photo-1518002171953-a080ee817e1f?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"adidas", "Adidas Ultraboost Light", "adidas-ultraboost-light", "running", false, 3299000}, "https://images.unsplash.com/photo-1587563871167-1ee9c731aefb?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"adidas", "Adidas Stan Smith", "adidas-stan-smith", "lifestyle", false, 1599000}, "https://images.unsplash.com/photo-1582588678413-dbf45f4823e9?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"vans", "Vans Old Skool", "vans-old-skool", "lifestyle", false, 899000}, "https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"vans", "Vans Authentic", "vans-authentic", "lifestyle", false, 799000}, "https://images.unsplash.com/photo-1560769629-975ec94e6a86?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"puma", "Puma Suede Classic", "puma-suede-classic", "lifestyle", false, 1299000}, "https://images.unsplash.com/photo-1608231387042-66d1773070a5?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"puma", "Puma RS-X", "puma-rs-x", "lifestyle", false, 1499000}, "https://images.unsplash.com/photo-1608231387042-66d1773070a5?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"new-balance", "New Balance 574", "new-balance-574", "lifestyle", false, 1399000}, "https://images.unsplash.com/photo-1539185441755-769473a23570?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"new-balance", "New Balance 990v6", "new-balance-990v6", "running", true, 3999000}, "https://images.unsplash.com/photo-1539185441755-769473a23570?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"asics", "Asics Gel-Kayano 30", "asics-gel-kayano-30", "running", false, 2799000}, "https://images.unsplash.com/photo-1584735935682-2f2b69dff9d2?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"jordan", "Air Jordan 1 Retro High OG", "air-jordan-1-retro-high-og", "lifestyle", true, 3499000}, "https://images.unsplash.com/photo-1552346154-21d32810aba3?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"jordan", "Air Jordan 4 Retro", "air-jordan-4-retro", "lifestyle", true, 4299000}, "https://images.unsplash.com/photo-1516478177764-9fe5bd7e9717?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"brodo", "Brodo Sella Oxford", "brodo-sella-oxford", "formal", false, 849000}, "https://images.unsplash.com/photo-1614252235316-8c857d38b5f4?auto=format&fit=crop&w=800&q=80"},
		{seedProduct{"brodo", "Brodo Lana Derby", "brodo-lana-derby", "formal", false, 949000}, "https://images.unsplash.com/photo-1614252235316-8c857d38b5f4?auto=format&fit=crop&w=800&q=80"},
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

		var productID uuid.UUID
		row := pool.QueryRow(ctx, `
			INSERT INTO products (id, brand_id, name, slug, category, is_limited)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, uuid.New(), brandID, p.name, p.slug, p.category, p.isLimited)
		if err := row.Scan(&productID); err != nil {
			log.Printf("skip product %s: %v", p.slug, err)
			continue
		}

		mustExecPool(ctx, pool, `
			UPDATE product_images SET url = $1 WHERE product_id = $2
		`, p.imageURL, productID)
		mustExecPool(ctx, pool, `
			INSERT INTO product_images (id, product_id, url, sort_order)
			SELECT $1, $2, $3, 0
			WHERE NOT EXISTS (SELECT 1 FROM product_images WHERE product_id = $2)
		`, uuid.New(), productID, p.imageURL)

		numOffers := 2 + rng.Intn(2)
		for i := 0; i < numOffers && i < len(storeList); i++ {
			storeID := storeList[(rng.Int())%len(storeList)]
			priceVariance := 1.0 + (rng.Float64()-0.5)*0.1
			price := p.basePrice * priceVariance
			mustExecPool(ctx, pool, `
				INSERT INTO product_offers (id, product_id, store_id, price, currency, in_stock, affiliate_url)
				VALUES ($1, $2, $3, $4, 'IDR', $5, $6)
				ON CONFLICT DO NOTHING
			`, uuid.New(), productID, storeID, price, rng.Intn(10) > 1, "https://kickpick.id/go/"+p.slug)
		}

		for d := 30; d >= 0; d-- {
			date := time.Now().AddDate(0, 0, -d)
			drift := 1.0 + (rng.Float64()-0.5)*0.08
			mustExecPool(ctx, pool, `
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
