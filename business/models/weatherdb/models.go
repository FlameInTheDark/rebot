// Code generated by sqlc. DO NOT EDIT.

package weatherdb

import ()

type Weather struct {
	ID        int32  `json:"id"`
	DiscordID string `json:"discord_id"`
	Location  string `json:"location"`
}
