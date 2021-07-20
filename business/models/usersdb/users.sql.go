// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package usersdb

import (
	"context"
)

const create = `-- name: Create :one
insert into users (discord_id)
values ($1) returning id, discord_id, created_at
`

func (q *Queries) Create(ctx context.Context, discordID string) (User, error) {
	row := q.db.QueryRowContext(ctx, create, discordID)
	var i User
	err := row.Scan(&i.ID, &i.DiscordID, &i.CreatedAt)
	return i, err
}

const find = `-- name: Find :one
select id, discord_id, created_at
from users
where discord_id = $1
`

func (q *Queries) Find(ctx context.Context, discordID string) (User, error) {
	row := q.db.QueryRowContext(ctx, find, discordID)
	var i User
	err := row.Scan(&i.ID, &i.DiscordID, &i.CreatedAt)
	return i, err
}