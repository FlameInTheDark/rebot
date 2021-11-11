package main

import (
	"context"
	"fmt"
	"github.com/FlameInTheDark/rebot/foundation/redisdb"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/app/commander/config"
	"github.com/FlameInTheDark/rebot/app/commander/service"
	"github.com/FlameInTheDark/rebot/foundation/database"
	"github.com/FlameInTheDark/rebot/foundation/discord"
	"github.com/FlameInTheDark/rebot/foundation/logs"
	"github.com/FlameInTheDark/rebot/foundation/queue"
)

func RunCommanderService(logger *zap.Logger) error {
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

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		logger.Error("database connection error", zap.Error(err))
		return err
	}

	sess, err := discord.NewDiscordSession(conf.Discord.Token)
	if err != nil {
		logger.Error("discord connection error", zap.Error(err))
		return err
	}

	rabbit, err := queue.NewRabbitmqConnection(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		conf.RabbitMQ.User,
		conf.RabbitMQ.Password,
		conf.RabbitMQ.Host,
		conf.RabbitMQ.Port,
	))

	rc, err := redisdb.NewConnection(conf.Redis.Host, conf.Redis.Port, conf.Redis.Password, conf.Redis.Database)
	if err != nil {
		logger.Error("redis connection error", zap.Error(err))
		return err
	}

	cmdr, err := service.NewCommander(db, rc, sess, rabbit, logger)
	if err != nil {
		logger.Error("worker creation error", zap.Error(err))
		return err
	}

	// Health check router
	r := chi.NewRouter()
	r.Use(
		logs.HttpLoggerMiddleware(logger),
		middleware.Recoverer,
	)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "OK")
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", conf.Http.Port), Handler: r}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			logger.Info("Service shutting down..")

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			err := srv.Shutdown(ctx)
			if err != nil {
				logger.Error("API server shutdown error", zap.Error(err))
			}
			select {
			case <-time.After(21 * time.Second):
				logger.Info("Not all connections done")
			case <-ctx.Done():

			}
		}
	}()

	return srv.ListenAndServe()
}
