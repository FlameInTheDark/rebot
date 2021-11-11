package service

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/service/discord/worker"
	"github.com/FlameInTheDark/rebot/business/service/users"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
)

type Commander struct {
	Users    *users.Service
	Discord  *worker.DiscordWorker
	Commands commandst.CommandsSender
	logger   *zap.Logger
}

func NewCommander(db *sqlx.DB, rc *redis.Client, sess *discordgo.Session, rabbit *amqp.Connection, logger *zap.Logger) (*Commander, error) {
	cmdService, err := commandst.NewRabbitCommandsTransport(rabbit, logger)
	if err != nil {
		return nil, err
	}

	return &Commander{
		Users:    users.NewUsersService(db),
		Discord:  worker.NewWorker(db, rc, sess, logger),
		Commands: cmdService,
		logger:   logger.With(zap.String("component", "discovery")),
	}, nil
}

func (c *Commander) Run() error {
	c.Discord.OnMessageHandler()
}
