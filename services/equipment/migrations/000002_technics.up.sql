create table if not exists equipment(
    id serial primary key,
    category_id int references categories(id) default 0,
    name varchar(256) not null,
    description text,
    additional_characteristics jsonb,
    created_at timestamp with time zone not null default current_timestamp,
    updated_at timestamp with time zone not null default current_timestamp
);

-- для быстрого поиска по категории
create index if not exists equipment_category_index on equipment(category_id);
