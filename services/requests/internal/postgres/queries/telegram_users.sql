-- name: SaveTelegramUser :one
INSERT INTO telegram_users (telegram_id, username, first_name, last_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (telegram_id) DO UPDATE
SET username = EXCLUDED.username,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    updated_at = NOW()
RETURNING *;

-- name: GetUserByTelegramID :one
SELECT *
FROM telegram_users
WHERE telegram_id = $1;

-- name: GetUserByID :one
SELECT *
FROM telegram_users
WHERE id = $1;
