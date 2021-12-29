package responses

import (
	"github.com/FlameInTheDark/rebot/business/models/usersdb"
	"time"
)

type UserResponse struct {
	ID        int32     `json:"id"`
	DiscordID string    `json:"discord_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(user *usersdb.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		DiscordID: user.DiscordID,
		CreatedAt: user.CreatedAt,
	}
}
