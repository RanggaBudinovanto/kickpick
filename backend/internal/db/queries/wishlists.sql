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

-- name: ListActiveAlertSubscribersForProduct :many
SELECT w.id AS wishlist_id, w.alert_type, u.id AS user_id, u.email, u.name
FROM wishlists w
JOIN users u ON u.id = w.user_id
WHERE w.product_id = $1 AND w.alert_active = true AND u.deleted_at IS NULL;
