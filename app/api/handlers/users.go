package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/FlameInTheDark/rebot/app/api/handlers/responses"
	"github.com/FlameInTheDark/rebot/business/service/users"
)

type UsersHandlersGroup struct {
	users *users.Service
}

func NewUsersHandlersGroup(users *users.Service) *UsersHandlersGroup {
	return &UsersHandlersGroup{users: users}
}

func (u *UsersHandlersGroup) GetUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.JSON(responses.Error(errors.New("id can not be empty")))
	}
	user, err := u.users.GetUser(c.Context(), id)
	if err != nil {
		return c.JSON(responses.Error(err))
	}
	return c.JSON(responses.NewUserResponse(user))
}
