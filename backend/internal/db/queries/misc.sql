-- name: GetSizeConversion :many
SELECT *
FROM size_conversion_matrix
WHERE reference_brand_id = $1 AND target_brand_id = $2 AND reference_size = $3;

-- name: GetUserSizePreferences :many
SELECT *
FROM user_size_preferences
WHERE user_id = $1;

-- name: UpsertUserSizePreference :one
INSERT INTO user_size_preferences (user_id, brand_id, size)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, brand_id) DO UPDATE SET size = EXCLUDED.size
RETURNING *;

-- name: GetLatestExchangeRate :one
SELECT *
FROM exchange_rates
WHERE base_currency = $1 AND target_currency = $2
ORDER BY recorded_date DESC
LIMIT 1;

-- name: SearchProductsAutocomplete :many
SELECT id, name, slug
FROM products
WHERE name ILIKE '%' || sqlc.arg('query')::text || '%'
ORDER BY name ASC
LIMIT 10;

-- name: SearchBrandsAutocomplete :many
SELECT id, name, slug
FROM brands
WHERE name ILIKE '%' || sqlc.arg('query')::text || '%'
ORDER BY name ASC
LIMIT 5;
