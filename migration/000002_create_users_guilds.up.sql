create table guilds
(
    id         bigint primary key,
    discord_id varchar   not null unique,
    created_at timestamp not null default now()
);