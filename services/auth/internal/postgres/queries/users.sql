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

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET username = COALESCE(NULLIF(sqlc.arg(username), ''), username),
    email = COALESCE(NULLIF(sqlc.arg(email), ''), email),
    updated_at = NOW()
WHERE id = $1
RETURNING id, username, email, password_hash, role_id, created_at, updated_at;
