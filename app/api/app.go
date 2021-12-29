package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/app/api/config"
	"github.com/FlameInTheDark/rebot/app/api/handlers"
	"github.com/FlameInTheDark/rebot/foundation/consul"
	"github.com/FlameInTheDark/rebot/foundation/database"
	"github.com/FlameInTheDark/rebot/foundation/logs"
)

// RunAPIServer create and start rest api server
func RunAPIServer(logger *zap.Logger) error {
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
		CertPath:   conf.Database.CetrPath,
		Logger:     logs.NewDBLogger(logger),
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		logger.Error("database connection error", zap.Error(err))
		return err
	}
	defer func() {
		dberr := db.Close()
		if dberr != nil {
			logger.Error("Database connection close error", zap.Error(dberr))
		}
	}()

	flog := logs.NewFiberLogger(logger)

	app := fiber.New()
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString(time.Now().String())
	})

	v1 := app.Group("/api/v1", flog.Middleware())

	handlers.API(v1, handlers.CreateServices(db), logger)

	consul, err := consul.NewConsulClient(conf.Consul.Address)
	if err != nil {
		logger.Error("Cannot create Consul client", zap.Error(err))
		return err
	}
	defer func() {
		cerr := consul.Close()
		if cerr != nil {
			logger.Error("Consul client close error", zap.Error(cerr))
		}
	}()
	err = consul.Register(conf.Consul.ServiceID.String(), conf.Consul.ServiceName, conf.Http.Port, nil)
	if err != nil {
		logger.Error("Cannot register service in consul", zap.Error(err))
		return err
	}
	defer func() {
		err := consul.Deregister(conf.Consul.ServiceID.String())
		if err != nil {
			logger.Error("Cannot deregister service", zap.Error(err))
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			logger.Info("Service shutting down..")

			ferr := app.Shutdown()
			if ferr != nil {
				logger.Error("API server shutdown error", zap.Error(ferr))
			}
		}
	}()
	return app.Listen(fmt.Sprintf(":%d", conf.Http.Port))
}
