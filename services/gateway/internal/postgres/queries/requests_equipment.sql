-- name: ListEquipmentForRequest :many
SELECT request_id, equipment_id, quantity, created_at, updated_at
FROM equipment_to_requests
WHERE request_id = $1;

