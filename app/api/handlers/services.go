package handlers

import (
	"github.com/jmoiron/sqlx"

	"github.com/FlameInTheDark/rebot/business/service/users"
)

// Services contains business logic
type Services struct {
	Users *users.Service
}

// CreateServices creates services
func CreateServices(db *sqlx.DB) *Services {
	var services Services

	services.Users = users.NewUsersService(db)

	return &services
}
