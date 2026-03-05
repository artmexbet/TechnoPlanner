-- name: AddEquipment :one
insert into equipment(category_id, name, description, additional_characteristics, quantity)
values($1, $2, $3, $4, $5) returning id;

-- name: DeleteEquipment :exec
delete from equipment where id = $1;

-- name: UpdateEquipment :exec
update equipment set category_id = $1, name = $2, description = $3, additional_characteristics = $4, quantity = $5 where id = $6;

-- name: GetEquipmentByID :one
select * from equipment where id = $1;

-- name: GetAllEquipment :many
select * from equipment order by id;

-- name: ReserveEquipment :exec
UPDATE equipment
SET reserved_quantity = reserved_quantity + $2
WHERE id = $1 AND (quantity - reserved_quantity) >= $2;

-- name: ReleaseEquipment :exec
UPDATE equipment
SET reserved_quantity = reserved_quantity - $2
WHERE id = $1 AND reserved_quantity >= $2;

-- name: GetAvailableQuantity :one
SELECT quantity - reserved_quantity AS available FROM equipment WHERE id = $1;
