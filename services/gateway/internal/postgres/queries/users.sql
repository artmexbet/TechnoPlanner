-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, role_id = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: ListPorters :many
SELECT * FROM users
WHERE role_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (username, email, role_id)
VALUES ($1, $2, $3)
RETURNING *;
