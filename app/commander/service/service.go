package service

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/service/discord/worker"
	"github.com/FlameInTheDark/rebot/business/service/discovery"
	"github.com/FlameInTheDark/rebot/business/service/users"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
	"github.com/FlameInTheDark/rebot/foundation/consul"
)

type Commander struct {
	Users     *users.Service
	Discord   *worker.DiscordWorker
	Registrar *RegistrarWorker
	logger    *zap.Logger
}

func NewCommander(db *sqlx.DB, rc *redis.Client, sess *discordgo.Session, cd *consul.Discovery, rabbit *amqp.Connection, logger *zap.Logger) (*Commander, error) {
	logger.Debug("Creating rabbit commands transport")
	cmdService, err := commandst.NewRabbitCommandsTransport(rabbit, logger)
	if err != nil {
		return nil, err
	}

	logger.Debug("Creating discovery service")
	cds := discovery.NewConsulDiscoveryService(cd, logger)

	logger.Debug("Creating commander service")
	return &Commander{
		Users:     users.NewUsersService(db),
		Discord:   worker.NewWorker(db, rc, sess, cmdService, logger),
		Registrar: NewRegistrarWorker(cds),
		logger:    logger.With(zap.String("component", "consul")),
	}, nil
}

func (c *Commander) Run() error {
	c.Registrar.AddRegistrarHandler(func(s consul.Service) {
		if data, ok := s.Meta["command_data"]; ok {
			c.logger.Debug("Registering command handler", zap.String("service-id", s.ID.String()), zap.Reflect("service-meta", s.Meta))
			meta, err := consul.ParseCommandMeta([]byte(data))
			if err != nil {
				return
			}
			for _, cmd := range *meta {
				c.Discord.AddCommandWorker(cmd.Command, cmd.Queue)
			}
		}
	})

	c.Registrar.Run("command")
	c.Discord.OnMessageHandler()
	return c.Discord.Open()
}
