-- name: ListProducts :many
SELECT
  p.*,
  b.name AS brand_name,
  b.slug AS brand_slug,
  COALESCE(MIN(po.price), 0)::numeric AS min_price,
  COALESCE(MAX(po.price), 0)::numeric AS max_price,
  COALESCE(AVG(r.rating), 0)::numeric AS avg_rating,
  COALESCE((SELECT url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.sort_order LIMIT 1), '')::text AS image_url
FROM products p
JOIN brands b ON b.id = p.brand_id
LEFT JOIN product_offers po ON po.product_id = p.id
LEFT JOIN reviews r ON r.product_id = p.id AND r.is_flagged = false
WHERE (sqlc.narg('brand_id')::uuid IS NULL OR p.brand_id = sqlc.narg('brand_id'))
  AND (sqlc.narg('brand_ids')::uuid[] IS NULL OR p.brand_id = ANY(sqlc.narg('brand_ids')::uuid[]))
  AND (sqlc.narg('category')::text IS NULL OR p.category = sqlc.narg('category'))
  AND (sqlc.narg('is_limited')::boolean IS NULL OR p.is_limited = sqlc.narg('is_limited'))
  AND (sqlc.narg('search')::text IS NULL OR p.name ILIKE '%' || sqlc.narg('search') || '%')
GROUP BY p.id, b.name, b.slug
HAVING (sqlc.narg('min_price')::numeric IS NULL OR (MIN(po.price) IS NOT NULL AND MIN(po.price) >= sqlc.narg('min_price')))
   AND (sqlc.narg('max_price')::numeric IS NULL OR (MIN(po.price) IS NOT NULL AND MIN(po.price) <= sqlc.narg('max_price')))
ORDER BY p.created_at DESC
LIMIT sqlc.arg('page_limit') OFFSET sqlc.arg('page_offset');

-- name: GetProductBySlug :one
SELECT
  p.*,
  b.name AS brand_name,
  b.slug AS brand_slug,
  COALESCE((SELECT url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.sort_order LIMIT 1), '')::text AS image_url
FROM products p
JOIN brands b ON b.id = p.brand_id
WHERE p.slug = $1;

-- name: ListOffersByProductID :many
SELECT po.*, s.name AS store_name, s.type AS store_type
FROM product_offers po
JOIN stores s ON s.id = po.store_id
WHERE po.product_id = $1
ORDER BY po.price ASC;

-- name: ListPriceHistoryByProductID :many
SELECT *
FROM price_history
WHERE product_id = $1
  AND recorded_date >= $2
ORDER BY recorded_date ASC;

-- name: RecordProductView :exec
INSERT INTO product_views (product_id)
VALUES ($1);

-- name: ListTrendingProducts :many
SELECT
  p.*,
  b.name AS brand_name,
  b.slug AS brand_slug,
  COALESCE(MIN(po.price), 0)::numeric AS min_price,
  COALESCE(MAX(po.price), 0)::numeric AS max_price,
  COALESCE(AVG(r.rating), 0)::numeric AS avg_rating,
  COALESCE((SELECT url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.sort_order LIMIT 1), '')::text AS image_url,
  count(pv.id) AS view_count
FROM products p
JOIN brands b ON b.id = p.brand_id
JOIN product_views pv ON pv.product_id = p.id AND pv.viewed_at >= now() - interval '7 days'
LEFT JOIN product_offers po ON po.product_id = p.id
LEFT JOIN reviews r ON r.product_id = p.id AND r.is_flagged = false
GROUP BY p.id, b.name, b.slug
ORDER BY view_count DESC, p.created_at DESC
LIMIT sqlc.arg('page_limit');

-- name: ListPriceDropProducts :many
WITH avg_prices AS (
  SELECT product_id, avg(price)::numeric AS avg_price_30d
  FROM price_history
  WHERE recorded_date >= CURRENT_DATE - interval '30 days'
  GROUP BY product_id
),
current_prices AS (
  SELECT product_id, min(price)::numeric AS current_min_price
  FROM product_offers
  WHERE in_stock = true
  GROUP BY product_id
)
SELECT
  p.*,
  b.name AS brand_name,
  b.slug AS brand_slug,
  cp.current_min_price AS min_price,
  cp.current_min_price AS max_price,
  COALESCE(AVG(r.rating), 0)::numeric AS avg_rating,
  COALESCE((SELECT url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.sort_order LIMIT 1), '')::text AS image_url,
  round((1 - (cp.current_min_price / ap.avg_price_30d)) * 100)::int AS drop_percent
FROM products p
JOIN brands b ON b.id = p.brand_id
JOIN avg_prices ap ON ap.product_id = p.id
JOIN current_prices cp ON cp.product_id = p.id
LEFT JOIN reviews r ON r.product_id = p.id AND r.is_flagged = false
WHERE cp.current_min_price < ap.avg_price_30d * 0.95
GROUP BY p.id, b.name, b.slug, cp.current_min_price, ap.avg_price_30d
ORDER BY drop_percent DESC
LIMIT sqlc.arg('page_limit');
