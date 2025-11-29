-- name: UpsertEquipmentCategories :exec
INSERT INTO equipment_category_links (equipment_id, category_id)
VALUES ($1, UNNEST($2::INT[]))
ON CONFLICT (equipment_id, category_id) DO NOTHING;

-- name: ClearEquipmentCategories :exec
DELETE FROM equipment_category_links
WHERE equipment_id = $1;

-- name: ListCategoriesForEquipment :many
SELECT c.*
FROM equipment_categories c
         JOIN equipment_category_links l ON c.id = l.category_id
WHERE l.equipment_id = $1 AND c.deleted_at IS NULL;

