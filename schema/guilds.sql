create table guilds
(
    id         bigint primary key,
    discord_id varchar   not null unique,
    created_at timestamp not null default now()
);

-- name: Find :one
select *
from guilds
where discord_id = $1;

-- name: Create :one
insert into guilds (discord_id)
values ($1) returning *;