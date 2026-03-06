-- name: GetAllCategories :many
select * from categories;

-- name: GetTechnicByCategory :many
select * from equipment where category_id = $1;

-- name: AddCategory :one
insert into categories(name, description) values($1, $2) returning id;

-- name: UpdateCategory :exec
update categories set name = $1, description = $2 where id = $3;

-- name: DeleteCategory :exec
delete from categories where id = $1;