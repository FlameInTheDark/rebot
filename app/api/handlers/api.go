package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/app/api/handlers/responses"
)

// API register api routes
func API(r fiber.Router, services *Services, logger *zap.Logger) {
	authGroup := NewAuthHandlersGroup()
	ag := r.Group("/auth")
	ag.Post("/callback", authGroup.OAuthCallbackHandler)

	usersGroup := NewUsersHandlersGroup(services.Users)
	ug := r.Group("/users")
	ug.Get("/:id", usersGroup.GetUserHandler)

	r.Get("/test", func(c *fiber.Ctx) error {
		err := services.Users.CreateUser(c.Context(), "160834320934764544")
		if err != nil {
			logger.Error("create user error", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(responses.Error(err))
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User has been created successfully"})
	})
}
