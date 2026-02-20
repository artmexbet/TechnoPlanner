-- name: SaveResponsible :one
INSERT INTO responsibles (id, username)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET username = EXCLUDED.username
RETURNING *;

-- name: GetResponsibleByID :one
SELECT *
FROM responsibles
WHERE id = $1;

-- name: GetResponsibleByUsername :one
SELECT *
FROM responsibles
WHERE username = $1;

-- name: ListResponsibles :many
SELECT *
FROM responsibles
ORDER BY username;

-- name: DeleteResponsible :exec
DELETE FROM responsibles
WHERE id = $1;
