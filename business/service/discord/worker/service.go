package worker

import "github.com/bwmarrin/discordgo"

type DiscordWorker struct {
	session *discordgo.Session
}

//NewWorker creates a new discord worker service
func NewWorker(session *discordgo.Session) *DiscordWorker {
	return &DiscordWorker{
		session: session,
	}
}

