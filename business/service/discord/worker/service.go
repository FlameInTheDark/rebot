package worker

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"

	"github.com/FlameInTheDark/rebot/business/service/guilds"
)

type DiscordWorker struct {
	session *discordgo.Session
	events  *EventRegistrar
	guilds  *guilds.GuildsService
	logger  *zap.Logger
}

//NewWorker creates a new discord worker service
func NewWorker(db *sqlx.DB, rc *redis.Client, session *discordgo.Session, logger *zap.Logger) *DiscordWorker {
	return &DiscordWorker{
		events:  NewEventRegistrar(),
		session: session,
		guilds:  guilds.NewGuildsService(db, rc, logger),
		logger:  logger,
	}
}

func (d *DiscordWorker) OnMessageHandler() {
	d.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if len(m.Content) < 2 {
			return
		}
		prefix, err := d.guilds.GetCommandPrefix(m.GuildID)
		if err != nil {
			d.logger.Error("cannot get prefix for guild", zap.String("discord-guild", m.GuildID))
			return
		}
		if string(m.Content[0]) != prefix {
			return
		}
		message := m.Content[1:]
		parts := strings.Split(message, " ")
		d.events.Notify(parts[0], MessageEvent{
			GuildID:  m.GuildID,
			UserID:   m.Author.ID,
			Username: m.Author.Username,
			Message:  message,
		})
	})
}
