// Code generated by sqlc. DO NOT EDIT.

package usersdb

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	DiscordID string    `json:"discord_id"`
	CreatedAt time.Time `json:"created_at"`
}
