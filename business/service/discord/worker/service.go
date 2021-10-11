package worker

import (
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

type DiscordWorker struct {
	session *discordgo.Session
	events  *EventRegistrar
	logger  *zap.Logger
}

//NewWorker creates a new discord worker service
func NewWorker(session *discordgo.Session) *DiscordWorker {
	return &DiscordWorker{
		events: NewEventRegistrar(),
		session: session,
	}
}
