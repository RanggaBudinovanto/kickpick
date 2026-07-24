-- name: GetOfferByID :one
SELECT *
FROM product_offers
WHERE id = $1;
