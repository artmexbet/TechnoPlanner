-- name: CreateRequest :one
INSERT INTO requests (telegram_user_id, request_text, schedule_time, address, end_time)
VALUES ($1, $2, $3, $4, $5)
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

-- name: GetRequestsByResponsibleID :many
SELECT *
FROM requests
WHERE porter_id = $1
ORDER BY created_at DESC;

-- name: AssignResponsible :one
UPDATE requests
SET porter_id = $2,
    status         = CASE
                         WHEN status = 'pending' THEN 'assigned'::request_status
                         ELSE status
        END,
    updated_at     = NOW()
WHERE id = $1
RETURNING *;

-- name: GetRequests :many
SELECT *
FROM requests
ORDER BY created_at DESC
OFFSET $1 LIMIT $2;
