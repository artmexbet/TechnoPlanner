-- name: GetAllCategories :many
select * from categories;

-- name: GetTechnicByCategory :many
select * from equipment where category_id = $1;

-- name: AddCategory :one
insert into categories(name) values($1) returning id;

-- name: UpdateCategoryName :exec
update categories set name = $1 where id = $2;

-- name: DeleteCategory :exec
delete from categories where id = $1;