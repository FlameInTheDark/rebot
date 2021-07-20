package users

import (
	"context"
	"github.com/FlameInTheDark/rebot/business/models/usersdb"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	users usersdb.Querier
}

func NewUsersService(db *sqlx.DB) *Service {
	return &Service{
		users: usersdb.New(db),
	}
}

func (s *Service) CreateUser(ctx context.Context, did string) error {
	_, err := s.users.Create(ctx, did)
	if err != nil {
		return err
	}
	return nil
}