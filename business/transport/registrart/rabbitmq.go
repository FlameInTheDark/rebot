package registrart

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var _ RegistrarSender = (*RabbitRegistrarTransport)(nil)
var _ RegistrarReceiver = (*RabbitRegistrarTransport)(nil)

type RabbitRegistrarTransport struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	logger *zap.Logger

	close chan struct{}
}

func NewRabbitRegistrarTransport(conn *amqp.Connection, logger *zap.Logger) (*RabbitRegistrarTransport, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create channel")
	}

	return &RabbitRegistrarTransport{
		conn:    conn,
		channel: ch,
		logger:  logger,
		close:   make(chan struct{}),
	}, nil
}

func (t *RabbitRegistrarTransport) Close() {
	close(t.close)
}

func (t *RabbitRegistrarTransport) getQueue(name string) (amqp.Queue, error) {
	return t.channel.QueueDeclare(name, true, true, false, false, nil)
}

func (t *RabbitRegistrarTransport) RegisterCommand(id uuid.UUID, command string) error {
	data, err := json.Marshal(&RegistrarMessage{
		ID:      id,
		Command: command,
	})
	if err != nil {
		return err
	}
	q, err := t.getQueue("service.registrar")
	if err != nil {
		return err
	}
	return t.channel.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func (t *RabbitRegistrarTransport) ReceiveRegisterRequests() (<-chan RegistrarMessage, error) {
	var ch = make(chan RegistrarMessage)

	q, err := t.getQueue("service.registrar")
	if err != nil {
		return nil, err
	}
	msgs, err := t.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func(rec <-chan amqp.Delivery, snd chan RegistrarMessage) {
		for {
			select {
			case <-t.close:
				return
			case msg := <-rec:
				var reg RegistrarMessage
				err := json.Unmarshal(msg.Body, &reg)
				if err != nil {
					t.logger.Error("command unmarshal error", zap.Error(err), zap.String("rabbit-queue", q.Name))
					continue
				}
				snd <- reg
			}
		}
	}(msgs, ch)

	return ch, nil
}
