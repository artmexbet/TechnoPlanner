-- name: CreateRawRequest :one
INSERT INTO raw_requests (telegram_id, username, first_name, last_name, raw_text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRawRequests :many
SELECT *
FROM raw_requests
WHERE ($1::text = '' OR status = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRawRequestByID :one
SELECT *
FROM raw_requests
WHERE id = $1;

-- name: MarkRawRequestProcessed :one
UPDATE raw_requests
SET status               = 'processed',
    processed_request_id = $2
WHERE id = $1
RETURNING *;

