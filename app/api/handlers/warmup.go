package handlers

import (
	"github.com/FlameInTheDark/rebot/business/service/users"
	"github.com/jmoiron/sqlx"
)

//Services contains business logic
type Services struct {
	Users *users.Service
}

//WarmupServices creates services
func WarmupServices(db *sqlx.DB) *Services {
	var services Services

	services.Users = users.NewUsersService(db)

	return &services
}
