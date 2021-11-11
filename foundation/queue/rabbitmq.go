package queue

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

func NewRabbitmqConnection(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, "rabbitmq connection error")
	}
	return conn, nil
}