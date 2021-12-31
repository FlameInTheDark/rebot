create table weather
(
    id         serial primary key,
    discord_id varchar not null unique,
    location   varchar not null default ''
);

-- name: Find :one
select location
from weather
where discord_id = $1;

-- name: Create :one
insert into weather (discord_id, location)
values ($1, $2) returning *;

-- name: Update :exec
update weather set location = $1 where discord_id = $1;

-- name: Insert :exec
insert into weather (discord_id, location)
values($1, $2)
    on conflict (discord_id)
do
update
    set location = $2;