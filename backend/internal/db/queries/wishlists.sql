-- name: ListWishlistByUser :many
SELECT w.*, p.name AS product_name, p.slug AS product_slug, p.brand_id
FROM wishlists w
JOIN products p ON p.id = w.product_id
WHERE w.user_id = $1
ORDER BY w.created_at DESC;

-- name: GetWishlistItem :one
SELECT *
FROM wishlists
WHERE id = $1;

-- name: AddWishlistItem :one
INSERT INTO wishlists (user_id, product_id, alert_active, alert_type)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, product_id) DO UPDATE SET alert_active = EXCLUDED.alert_active
RETURNING *;

-- name: RemoveWishlistItem :exec
DELETE FROM wishlists
WHERE id = $1 AND user_id = $2;

-- name: SetWishlistAlert :one
UPDATE wishlists
SET alert_active = $3, alert_type = $4
WHERE id = $1 AND user_id = $2
RETURNING *;
