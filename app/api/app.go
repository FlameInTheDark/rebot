package main

import (
	"context"
	"fmt"
	"github.com/FlameInTheDark/rebot/app/api/handlers"
	"github.com/FlameInTheDark/rebot/foundation/database"
	"github.com/FlameInTheDark/rebot/foundation/logs"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunAPIServer(logger *zap.Logger) error {
	conf, err := LoadConfig()
	if err != nil {
		logger.Error("configuration not loaded", zap.Error(err))
		return err
	}

	dbConfig := database.Config{
		Host:       conf.Database.Host,
		Port:       conf.Database.Port,
		Database:   conf.Database.Database,
		Username:   conf.Database.Username,
		Password:   conf.Database.Password,
		DisableTLS: conf.Database.DisableTLS,
		CertPath:   conf.Database.CetrPath,
		Logger: logs.NewDBLogger(logger),
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		logger.Error("database connection error", zap.Error(err))
		return err
	}

	r := chi.NewRouter()
	r.Use(logs.HttpLoggerMiddleware(logger))

	handlers.API(r, handlers.WarmupServices(db), logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", conf.Http.Port), Handler: r}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			// sig is a ^C, handle it
			logger.Info("Service shutting down..")

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// start http shutdown
			srv.Shutdown(ctx)

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				logger.Info("Not all connections done")
			case <-ctx.Done():

			}
		}
	}()
	return srv.ListenAndServe()
}
