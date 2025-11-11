-- name: CreateRequest :one
INSERT INTO requests (telegram_user_id, request_text)
VALUES ($1, $2)
RETURNING *;

-- name: GetRequestsByTelegramUserID :many
SELECT *
FROM requests
WHERE telegram_user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRequestStatus :one
UPDATE requests
SET status     = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetRequestByID :one
SELECT *
FROM requests
WHERE id = $1;