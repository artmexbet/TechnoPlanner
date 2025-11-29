-- name: CreateEquipmentCategory :one
INSERT INTO equipment_categories (name, description, created_by, updated_by)
VALUES ($1, $2, $3, $3)
RETURNING *;

-- name: UpdateEquipmentCategory :one
UPDATE equipment_categories
SET name = $2,
    description = $3,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = $4
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteEquipmentCategory :exec
UPDATE equipment_categories
SET deleted_at = CURRENT_TIMESTAMP,
    updated_by = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEquipmentCategories :many
SELECT *
FROM equipment_categories
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

