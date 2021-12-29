create table guilds
(
    id             serial primary key,
    discord_id     varchar   not null unique,
    command_prefix varchar   not null default '!',
    created_at     timestamp not null default now()
);