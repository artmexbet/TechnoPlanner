create table if not exists technics(
    id uuid primary key default(uuid_generate_v4()),
    category_id uuid references categories(id) default '00000000-0000-0000-0000-000000000000',
    name varchar(256) not null,
    description text,
    additional_characteristics jsonb,
    created_at timestamp with time zone not null default current_timestamp,
    updated_at timestamp with time zone not null default current_timestamp
);

-- я думаю он не нужен...
create index if not exists technics_category_index on technics(category);