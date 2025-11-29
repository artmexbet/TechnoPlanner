-- name: AssignResponsible :one
UPDATE requests
SET responsible_user_id = $2,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ListRequests :many
SELECT *
FROM requests
WHERE deleted_at IS NULL
  AND ($1::UUID = '00000000-0000-0000-0000-000000000000' OR responsible_user_id = $1)
ORDER BY created_at DESC;

-- name: GetRequestByID :one
SELECT *
FROM requests
WHERE id = $1 AND deleted_at IS NULL;
