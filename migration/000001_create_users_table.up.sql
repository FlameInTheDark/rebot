create table users
(
    id         serial primary key,
    discord_id varchar   not null unique,
    created_at timestamp not null default now()
);