-- name: AddTechnic :batchone
INSERT INTO technics (name, description, quantity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: BatchGetTechnicsByRequestID :batchmany
SELECT t.* FROM technics_to_requests tr
JOIN technics t ON tr.technic_id = t.id AND t.quantity > 0
WHERE tr.request_id = $1
ORDER BY t.created_at DESC;

-- name: GetTechnicsByRequestID :many
SELECT t.* FROM technics_to_requests tr
JOIN technics t ON tr.technic_id = t.id AND t.quantity > 0
WHERE tr.request_id = $1
ORDER BY t.created_at DESC;
