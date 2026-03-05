-- name: AddEquipment :batchone
INSERT INTO equipment (name, description, quantity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: BatchGetEquipmentByRequestID :batchmany
SELECT t.* FROM equipment_to_requests tr
JOIN equipment t ON tr.equipment_id = t.id AND t.quantity > 0
WHERE tr.request_id = $1
ORDER BY t.created_at DESC;

-- name: GetEquipmentByRequestID :many
SELECT t.* FROM equipment_to_requests tr
JOIN equipment t ON tr.equipment_id = t.id AND t.quantity > 0
WHERE tr.request_id = $1
ORDER BY t.created_at DESC;

-- name: UpsertEquipment :exec
INSERT INTO equipment (id, name, description, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
    SET name        = EXCLUDED.name,
        description = EXCLUDED.description,
        quantity    = EXCLUDED.quantity,
        updated_at  = CURRENT_TIMESTAMP;

-- name: DeleteEquipmentByID :exec
DELETE FROM equipment WHERE id = $1;

