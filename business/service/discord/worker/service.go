package worker

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/service/guilds"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
)

type DiscordWorker struct {
	session       *discordgo.Session
	events        *EventRegistrar
	guilds        *guilds.GuildsService
	commandSender commandst.CommandsSender
	logger        *zap.Logger
}

//NewWorker creates a new discord worker service
func NewWorker(
	db *sqlx.DB,
	rc *redis.Client,
	session *discordgo.Session,
	cmdst commandst.CommandsSender,
	logger *zap.Logger,
) *DiscordWorker {
	return &DiscordWorker{
		events:        NewEventRegistrar(logger),
		session:       session,
		guilds:        guilds.NewGuildsService(db, rc, logger),
		commandSender: cmdst,
		logger:        logger,
	}
}

func (d *DiscordWorker) Open() error {
	return d.session.Open()
}

func (d *DiscordWorker) AddCommandWorker(command, queue string) {
	d.events.Register(command, NewRabbitCommandObserver(queue, d.commandSender, d.logger))
}

//OnMessageHandler ...
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
		if !strings.HasPrefix(m.Content, prefix) {
			return
		}
		message := m.Content[1:]
		parts := strings.Split(message, " ")
		if _, ok := d.events.events[parts[0]]; !ok {
			return
		}
		d.logger.Debug("Triggered command", zap.String("command", parts[0]))
		d.events.Notify(strings.ToLower(parts[0]), MessageEvent{
			GuildID:   m.GuildID,
			UserID:    m.Author.ID,
			ChannelID: m.ChannelID,
			Username:  m.Author.Username,
			Message:   strings.Join(parts[1:], " "),
		})
	})
}
