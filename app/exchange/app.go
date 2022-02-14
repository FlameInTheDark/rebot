package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/app/weather/config"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
	"github.com/FlameInTheDark/rebot/foundation/consul"
	"github.com/FlameInTheDark/rebot/foundation/database"
	"github.com/FlameInTheDark/rebot/foundation/discord"
	"github.com/FlameInTheDark/rebot/foundation/logs"
	"github.com/FlameInTheDark/rebot/foundation/queue"
	"github.com/FlameInTheDark/rebot/foundation/redisdb"
)

// RunExchangeService runs an exchange service
func RunExchangeService(logger *zap.Logger) error {
	conf, err := config.GetConfig()
	if err != nil {
		logger.Error("configuration not loaded", zap.Error(err))
		return err
	}

	logger.Debug("ENVs loaded", zap.Reflect("config", conf))

	dbConfig := database.Config{
		Host:       conf.Database.Host,
		Port:       conf.Database.Port,
		Database:   conf.Database.Database,
		Username:   conf.Database.Username,
		Password:   conf.Database.Password,
		DisableTLS: conf.Database.DisableTLS,
		CertPath:   conf.Database.CertPath,
		Logger:     logs.NewDBLogger(logger),
	}

	logger.Debug("Creating database connection")
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		logger.Error("database connection error", zap.Error(err))
		return err
	}
	defer func() {
		logger.Info("Closing database connection")
		derr := db.Close()
		if derr != nil {
			logger.Error("Database connection close error", zap.Error(derr))
		}
	}()

	logger.Debug("Creating Discord session")
	sess, err := discord.NewDiscordSession(conf.Discord.Token)
	if err != nil {
		logger.Error("discord connection error", zap.Error(err))
		return err
	}
	defer func() {
		logger.Info("Closing discord session")
		serr := sess.Close()
		if serr != nil {
			logger.Error("Discord session close error", zap.Error(serr))
		}
	}()

	logger.Debug("Creating RabbitMQ connection")
	rabbit, err := queue.NewRabbitmqConnection(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		conf.RabbitMQ.User,
		conf.RabbitMQ.Password,
		conf.RabbitMQ.Host,
		conf.RabbitMQ.Port,
	))
	if err != nil {
		logger.Error("Rabbit connection error", zap.Error(err))
		return err
	}
	defer func() {
		logger.Info("Closing rabbit connection")
		rerr := rabbit.Close()
		if rerr != nil {
			logger.Error("Rabbit connection close error", zap.Error(rerr))
		}
	}()

	logger.Debug("Creating Redis connection")
	rc, err := redisdb.NewConnection(conf.Redis.Host, conf.Redis.Port, conf.Redis.Password, conf.Redis.Database)
	if err != nil {
		logger.Error("redis connection error", zap.Error(err))
		return err
	}
	defer func() {
		logger.Info("Closing redis client")
		rcerr := rc.Close()
		if rcerr != nil {
			logger.Error("Redis client close error", zap.Error(rcerr))
		}
	}()

	rbc, err := commandst.NewRabbitCommandsTransport(rabbit, logger)
	if err != nil {
		logger.Error("Rabbit command transport creating error", zap.Error(err))
		return err
	}

	logger.Debug("Creating Consul client")
	cd, err := consul.NewConsulClient(conf.Consul.Address)
	if err != nil {
		logger.Error("Cannot create Consul client", zap.Error(err))
		return err
	}

	app := fiber.New()
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString(time.Now().String())
	})

	meta, err := consul.MarshalCommandMeta(consul.CommandMetaInfo{{"w", "weather"}, {"ww", "wweather"}})
	if err != nil {
		logger.Error("Consul meta parsing error", zap.Error(err))
		return err
	}

	logger.Debug("Registering service", zap.String("service-name", conf.Consul.ServiceName))
	err = cd.Register(conf.Consul.ServiceID.String(), conf.Consul.ServiceName, conf.HTTP.Port, map[string]string{"command_data": meta})
	if err != nil {
		logger.Error("Cannot register service in consul", zap.Error(err))
		return err
	}
	defer func() {
		err := cd.Deregister(conf.Consul.ServiceID.String())
		if err != nil {
			logger.Error("Cannot deregister service", zap.Error(err))
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			logger.Info("Service shutting down..")

			logger.Info("Shutting down http endpoint")
			aerr := app.Shutdown()
			if aerr != nil {
				logger.Error("API server shutdown error", zap.Error(aerr))
			}
		}
	}()

	logger.Debug("Listening API")
	return app.Listen(fmt.Sprintf(":%d", conf.HTTP.Port))
}
