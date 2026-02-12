-- name: AddEquipment :one
insert into equipment(category_id, name, description, additional_characteristics) values($1, $2, $3, $4) returning id;

-- name: DeleteEquipment :exec
delete from equipment where id = $1;

-- name: UpdateEquipment :exec
update equipment set category_id = $1, name = $2, description = $3, additional_characteristics = $4 where id = $5;

-- name: GetEquipmentByID :one
select * from equipment where id = $1;