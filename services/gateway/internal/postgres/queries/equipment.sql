-- name: CreateEquipment :one
INSERT INTO equipment (name, description, quantity, created_by, updated_by)
VALUES ($1, $2, $3, $4, $4)
RETURNING *;

-- name: UpdateEquipment :one
UPDATE equipment
SET name = $2,
    description = $3,
    quantity = $4,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = $5
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteEquipment :exec
UPDATE equipment
SET deleted_at = CURRENT_TIMESTAMP,
    updated_by = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetEquipmentByID :one
SELECT *
FROM equipment
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEquipment :many
SELECT *
FROM equipment
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

