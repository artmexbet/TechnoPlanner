CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists categories(
    id serial primary key,
    name varchar(256) not null unique,
    description text,
    created_at timestamp with time zone not null default current_timestamp,
    updated_at timestamp with time zone not null default current_timestamp
);

-- дефолтная категория для оборудования без категории
insert into categories (id, name, description) values (0, 'Нет категории', 'Категория по умолчанию') on conflict do nothing;
