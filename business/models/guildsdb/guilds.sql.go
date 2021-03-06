// Code generated by sqlc. DO NOT EDIT.
// source: guilds.sql

package guildsdb

import (
	"context"
)

const create = `-- name: Create :one
insert into guilds (discord_id)
values ($1) returning id, discord_id, command_prefix, created_at
`

func (q *Queries) Create(ctx context.Context, discordID string) (Guild, error) {
	row := q.db.QueryRowContext(ctx, create, discordID)
	var i Guild
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.CommandPrefix,
		&i.CreatedAt,
	)
	return i, err
}

const find = `-- name: Find :one
select id, discord_id, command_prefix, created_at
from guilds
where discord_id = $1
`

func (q *Queries) Find(ctx context.Context, discordID string) (Guild, error) {
	row := q.db.QueryRowContext(ctx, find, discordID)
	var i Guild
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.CommandPrefix,
		&i.CreatedAt,
	)
	return i, err
}
