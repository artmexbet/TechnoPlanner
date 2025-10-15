-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: FindUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, password_hash = $4, role_id = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;