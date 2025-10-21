-- name: AddPlace :one
INSERT INTO places (name, description, latitude, longitude)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPlaces :many
SELECT * FROM places
ORDER BY created_at DESC;
