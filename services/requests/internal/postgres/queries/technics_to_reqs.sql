-- name: AssignTechnicToRequest :batchone
INSERT INTO technics_to_requests (technic_id, request_id, quantity)
VALUES ($1, $2, $3)
RETURNING *;