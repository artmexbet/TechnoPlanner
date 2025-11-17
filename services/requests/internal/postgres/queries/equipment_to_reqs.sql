-- name: AssignEquipmentToRequest :batchone
INSERT INTO equipment_to_requests (equipment_id, request_id, quantity)
VALUES ($1, $2, $3)
RETURNING *;