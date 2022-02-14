package exchange

import (
	"bytes"

	"github.com/jmoiron/sqlx"

	"github.com/FlameInTheDark/rebot/business/models/exchangedb"
)

type Service struct {
	db exchangedb.Querier
}

func NewService(db sqlx.DB) *Service {
	return &Service{db: exchangedb.New(db)}
}

func (s *Service) GetRates(base string, symbols []string, amount float64) (*bytes.Buffer, error) {

}
