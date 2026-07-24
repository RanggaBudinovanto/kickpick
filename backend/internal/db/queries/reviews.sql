-- name: ListReviewsByProduct :many
SELECT r.*, u.name AS user_name
FROM reviews r
JOIN users u ON u.id = r.user_id
WHERE r.product_id = $1 AND r.is_flagged = false
ORDER BY r.created_at DESC;

-- name: CreateReview :one
INSERT INTO reviews (product_id, user_id, rating, comment, fit_feedback)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ReportReview :exec
INSERT INTO review_reports (review_id, reported_by, reason)
VALUES ($1, $2, $3)
ON CONFLICT (review_id, reported_by) DO NOTHING;

-- name: FlagReviewIfReportsExceedThreshold :exec
UPDATE reviews r
SET is_flagged = true
WHERE r.id = $1
  AND (SELECT count(*) FROM review_reports rr WHERE rr.review_id = r.id) >= 3;
