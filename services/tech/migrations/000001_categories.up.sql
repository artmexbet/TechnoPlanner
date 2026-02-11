CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists categories(
    id uuid primary key default uuid_generate_v4(),
    name varchar(256) not null unique
);

insert into categories values('00000000-0000-0000-000000000000', 'Нет категории');