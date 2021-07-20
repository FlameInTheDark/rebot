create table users
(
    id         bigint primary key,
    discord_id varchar   not null unique,
    created_at timestamp not null default now()
);