package registrart

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
)

var _ RegistrarSender = (*RabbitRegistrarTransport)(nil)
var _ RegistrarReceiver = (*RabbitRegistrarTransport)(nil)

type RabbitRegistrarTransport struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	rw       sync.RWMutex
	handlers []RegistrarHandler

	logger *zap.Logger

	close chan struct{}
}

func (t *RabbitRegistrarTransport) Start() error {
	q, err := t.getQueue("service.registrar")
	if err != nil {
		return err
	}
	msgs, err := t.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(msgs <-chan amqp.Delivery, close chan struct{}) {
		for {
			select {
			case msg := <-msgs:
				var reg RegistrarMessage
				err := json.Unmarshal(msg.Body, &reg)
				if err != nil {
					t.logger.Error("command unmarshal error", zap.Error(err), zap.String("rabbit-queue", q.Name))
					continue
				}

				t.rw.RLock()
				for _, h := range t.handlers {
					h(reg.ID, reg.Command)
				}
				t.rw.RUnlock()
			case <-close:
				return
			}
		}
	}(msgs, t.close)

	return nil
}

//AddHandler adds a handler to the worker
func (t *RabbitRegistrarTransport) AddHandler(handler RegistrarHandler) {
	t.rw.Lock()
	t.handlers = append(t.handlers, handler)
	t.rw.Unlock()
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
