-- name: ListBrands :many
SELECT *
FROM brands
ORDER BY name ASC;

-- name: GetBrandBySlug :one
SELECT *
FROM brands
WHERE slug = $1;
