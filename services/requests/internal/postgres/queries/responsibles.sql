-- name: SaveResponsible :one
INSERT INTO porters (id, username)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET username = EXCLUDED.username
RETURNING *;

-- name: GetResponsibleByID :one
SELECT *
FROM porters
WHERE id = $1;

-- name: GetResponsibleByUsername :one
SELECT *
FROM porters
WHERE username = $1;

-- name: ListResponsibles :many
SELECT *
FROM porters
ORDER BY username;

-- name: DeleteResponsible :exec
DELETE FROM porters
WHERE id = $1;
