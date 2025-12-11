-- name: AddTechnic :one
insert into technics(category_id, name, description, additional_characteristics) values($1, $2, $3, $4) returning id;

-- name: DeleteTechnic :exec
delete from technics where id = $1;

-- name: UpdateTechnic :exec
update technics set category_id = $1, name = $2, description = $3, additional_characteristics = $4 where id = $5;

-- name: GetTechnicByID :one
select * from technics where id = $1;