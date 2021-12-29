package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//NewDiscordSession create new discord session with token
func NewDiscordSession(token string) (*discordgo.Session, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, errors.Wrap(err, "discord session creation failed")
	}

	sess.Identify.Intents = discordgo.IntentsAll

	return sess, nil
}
