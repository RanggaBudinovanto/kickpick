-- name: UpsertProduct :one
INSERT INTO products (id, brand_id, name, slug, category, is_limited)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
ON CONFLICT (slug) DO UPDATE SET
  name = EXCLUDED.name,
  category = EXCLUDED.category,
  is_limited = EXCLUDED.is_limited,
  updated_at = now()
RETURNING *;

-- name: ReplaceProductImage :exec
WITH deleted AS (
  DELETE FROM product_images WHERE product_id = $1
)
INSERT INTO product_images (id, product_id, url, sort_order)
SELECT gen_random_uuid(), $1, $2, 0
WHERE $2 != '';

-- name: GetStoreByName :one
SELECT * FROM stores WHERE name = $1;

-- name: CreateStore :one
INSERT INTO stores (id, name, type, affiliate_network)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetOfferByProductAndStore :one
SELECT * FROM product_offers WHERE product_id = $1 AND store_id = $2;

-- name: CreateProductOffer :one
INSERT INTO product_offers (id, product_id, store_id, price, currency, in_stock, affiliate_url)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProductOffer :exec
UPDATE product_offers
SET price = $3, currency = $4, in_stock = $5, affiliate_url = $6, scraped_at = now()
WHERE product_id = $1 AND store_id = $2;

-- name: UpsertPriceHistoryToday :exec
INSERT INTO price_history (id, product_id, store_id, price, recorded_date)
VALUES (gen_random_uuid(), $1, $2, $3, CURRENT_DATE)
ON CONFLICT (product_id, store_id, recorded_date) DO UPDATE SET price = EXCLUDED.price;
