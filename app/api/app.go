package main

import (
	"fmt"
	"github.com/FlameInTheDark/rebot/foundation/logs"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func API(logger *zap.Logger) (*http.Server, error) {
	conf, err := LoadConfig()
	if err != nil {
		logger.Error("configuration not loaded", zap.Error(err))
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(logs.HttpLoggerMiddleware(logger))

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"message":"hello"}`))
		if err != nil {
			logger.Error("Response error", zap.Error(err))
		}
	})

	return &http.Server{Addr: fmt.Sprintf(":%d", conf.Http.Port), Handler: r}, nil
}
