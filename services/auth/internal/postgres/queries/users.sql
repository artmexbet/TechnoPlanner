-- name: FindUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;