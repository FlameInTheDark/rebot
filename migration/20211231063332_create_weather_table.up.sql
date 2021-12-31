create table weather
(
    id         serial primary key,
    discord_id varchar not null unique,
    location   varchar not null default ''
);