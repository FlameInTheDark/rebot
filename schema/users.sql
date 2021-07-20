create table users
(
    id         serial primary key,
    discord_id varchar   not null unique,
    created_at timestamp not null default now()
);

-- name: Find :one
select *
from users
where discord_id = $1;

-- name: Create :one
insert into users (discord_id)
values ($1) returning *;