create table exchange
(
    id         serial primary key,
    discord_id varchar not null unique,
    base       varchar not null default 'USD'
);

-- name: Find :one
select base
from exchange
where discord_id = $1;

-- name: Create :one
insert into exchange (discord_id, base)
values ($1, $2) returning *;

-- name: Update :exec
update exchange
set base = $1
where discord_id = $1;

-- name: Insert :exec
insert into exchange (discord_id, base)
values ($1, $2) on conflict (discord_id)
do
update
    set base = $2;