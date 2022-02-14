package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/FlameInTheDark/rebot/app/api/handlers/requests"
	"github.com/FlameInTheDark/rebot/app/api/handlers/responses"
	"github.com/FlameInTheDark/rebot/foundation/validation"
)

// AuthHandlersGroup contains handlers for the auth API
type AuthHandlersGroup struct {
}

// NewAuthHandlersGroup creates a new AuthHandlersGroup
func NewAuthHandlersGroup() *AuthHandlersGroup {
	return &AuthHandlersGroup{}
}

// OAuthCallbackHandler accept discord oauth requests
func (a *AuthHandlersGroup) OAuthCallbackHandler(c *fiber.Ctx) error {
	var req requests.OauthCode
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Error(err))
	}

	verr := validation.ValidateStruct(req)
	if verr != nil {
		return c.JSON(verr)
	}

	return c.Status(fiber.StatusOK).JSON(req)
}
