-- name: AddPlace :one
INSERT INTO places (name, description, latitude, longitude)
VALUES ($1, $2, $3, $4)
RETURNING *;