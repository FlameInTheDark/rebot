package users

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/FlameInTheDark/rebot/business/models/usersdb"
)

type Service struct {
	users usersdb.Querier
}

func NewUsersService(db *sqlx.DB) *Service {
	return &Service{
		users: usersdb.New(db),
	}
}

//CreateUser creates new user in database with specified discord id
func (s *Service) CreateUser(ctx context.Context, did string) error {
	_, err := s.users.Create(ctx, did)
	if err != nil {
		return err
	}
	return nil
}

//GetUser returns user by discord id
func (s *Service) GetUser(ctx context.Context, did string) (*usersdb.User, error) {
	u, err := s.users.Find(ctx, did)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
