-- name: CreateUser :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: SetUserEmailVerified :exec
UPDATE users
SET email_verified = true, updated_at = now()
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET name = $2,
    onboarding_focus = $3,
    preferred_language = $4,
    preferred_currency = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;
