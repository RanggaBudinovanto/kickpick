-- name: ListNotificationsByUser :many
SELECT *
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUnreadNotifications :one
SELECT count(*)
FROM notifications
WHERE user_id = $1 AND is_read = false;

-- name: MarkNotificationRead :exec
UPDATE notifications
SET is_read = true
WHERE id = $1 AND user_id = $2;

-- name: CreateNotification :one
INSERT INTO notifications (user_id, type, title, body, action_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
