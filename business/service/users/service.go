package users

import (
	"github.com/FlameInTheDark/rebot/business/models/usersdb"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	users usersdb.Querier
}

func NewUsersService(db sqlx.DB) *Service {
	return &Service{
		users: usersdb.New(db),
	}
}

