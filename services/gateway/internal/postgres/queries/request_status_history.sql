-- name: InsertRequestStatusHistory :one
INSERT INTO request_status_history (request_id, status, comment, changed_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListRequestStatusHistory :many
SELECT *
FROM request_status_history
WHERE request_id = $1
ORDER BY changed_at DESC;

