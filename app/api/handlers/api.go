package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//API register api routes
func API(r chi.Router, services *Services, logger *zap.Logger) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authGroup := NewAuthHandlersGroup()
			r.Post("/callback", authGroup.OAuthCallbackHandler)
		})

		r.Group(func(r chi.Router) {

		})
	})

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		err := services.Users.CreateUser(r.Context(), "160834320934764544")
		if err != nil {
			logger.Error("create user error", zap.Error(err))
		}
	})
}
