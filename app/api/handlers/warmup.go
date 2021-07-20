package handlers

import (
	"github.com/FlameInTheDark/rebot/business/service/users"
	"github.com/jmoiron/sqlx"
)

type Services struct {
	Users *users.Service
}

func WarmupServices(db *sqlx.DB) *Services {
	var services Services

	services.Users = users.NewUsersService(db)

	return &services
}
