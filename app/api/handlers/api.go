package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func API(r chi.Router, services *Services, logger *zap.Logger) {
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		err := services.Users.CreateUser(r.Context(), "160834320934764544")
		if err != nil {
			logger.Error("create user error", zap.Error(err))
		}
	})

}
