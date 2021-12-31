package commandst

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var _ CommandsSender = (*RabbitCommandsTransport)(nil)
var _ CommandsReceiver = (*RabbitCommandsTransport)(nil)

type RabbitCommandsTransport struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	rw       sync.RWMutex
	handlers map[string]ReceiverHandler

	logger *zap.Logger

	close chan struct{}
}

func (t *RabbitCommandsTransport) AddHandler(command string, handler ReceiverHandler) {
	t.rw.Lock()
	defer t.rw.Unlock()
	if _, ok := t.handlers[command]; !ok {
		t.handlers[command] = handler
	}
}

func (t *RabbitCommandsTransport) Start(command string) error {
	var th ReceiverHandler
	if _, ok := t.handlers[command]; ok {
		th = t.handlers[command]
	} else {
		return errors.New("no commands found")
	}
	q, err := t.getQueue(fmt.Sprintf("commandst.%s", command))
	if err != nil {
		return err
	}

	msgs, err := t.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(msgs <-chan amqp.Delivery, handler ReceiverHandler, close chan struct{}) {
		for {
			select {
			case msg := <-msgs:
				var cmd CommandMessage
				err := json.Unmarshal(msg.Body, &cmd)
				if err != nil {
					t.logger.Error("command unmarshal error", zap.Error(err), zap.String("rabbit-queue", q.Name))
					continue
				}
				t.logger.Debug("Got message", zap.Reflect("rabbit-message", cmd))
				handler(cmd)
			case <-close:
				return
			}
		}
	}(msgs, th, t.close)

	t.logger.Debug("Started listening command channel", zap.String("rabbit-queue", fmt.Sprintf("commandst.%s", command)))

	return nil
}

func NewRabbitCommandsTransport(conn *amqp.Connection, logger *zap.Logger) (*RabbitCommandsTransport, error) {
	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Cannot create rabbit channel", zap.Error(err))
		return nil, errors.Wrap(err, "cannot create channel")
	}

	return &RabbitCommandsTransport{
		conn:    conn,
		handlers: make(map[string]ReceiverHandler),
		channel: ch,
		logger:  logger,
		close:   make(chan struct{}),
	}, nil
}

func (t *RabbitCommandsTransport) Close() {
	close(t.close)
}

func (t *RabbitCommandsTransport) getQueue(name string) (amqp.Queue, error) {
	return t.channel.QueueDeclare(name, true, true, false, false, nil)
}

func (t *RabbitCommandsTransport) SendCommand(cmd CommandMessage, queue string) error {
	data, err := json.Marshal(&cmd)
	if err != nil {
		return err
	}
	q, err := t.getQueue(fmt.Sprintf("commandst.%s", queue))
	if err != nil {
		return err
	}
	return t.channel.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func (t *RabbitCommandsTransport) ReceiveCommands(command string) (<-chan CommandMessage, error) {
	var ch = make(chan CommandMessage)

	q, err := t.getQueue(fmt.Sprintf("commandst.%s", command))
	if err != nil {
		return nil, err
	}
	msgs, err := t.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func(rec <-chan amqp.Delivery, snd chan CommandMessage) {
		for {
			select {
			case <-t.close:
				return
			case msg := <-rec:
				var command CommandMessage
				err := json.Unmarshal(msg.Body, &command)
				if err != nil {
					t.logger.Error("command unmarshal error", zap.Error(err), zap.String("rabbit-queue", q.Name))
					continue
				}
				snd <- command
			}
		}
	}(msgs, ch)

	return ch, nil
}

func (t *RabbitCommandsTransport) ReceivePings(command string) error {
	q, err := t.getQueue(fmt.Sprintf("commandst.%s.ping", command))
	if err != nil {
		return err
	}
	msgs, err := t.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(rec <-chan amqp.Delivery, ch *amqp.Channel, command string) {
		for {
			select {
			case <-t.close:
				return
			case msg := <-rec:
				var ping PingMessage
				err := json.Unmarshal(msg.Body, &ping)
				if err != nil {
					t.logger.Error("command unmarshal error", zap.Error(err), zap.String("rabbit-queue", fmt.Sprintf("commandst.%s.ping", command)))
					continue
				}
				if ping.Status == RCommandPing {
					pq, err := ch.QueueDeclare(
						fmt.Sprintf("commandst.%s.pong", command),
						true,
						true,
						false,
						false,
						nil,
					)
					if err != nil {
						t.logger.Error("pong queue declare error", zap.Error(err), zap.String("rabbit-queue", fmt.Sprintf("commandst.%s.pong", command)))

						continue
					}
					ping.Status = RCommandPong
					data, err := json.Marshal(&ping)
					if err != nil {
						t.logger.Error("pong marshal error", zap.Error(err), zap.String("rabbit-queue", fmt.Sprintf("commandst.%s.pong", command)))
						continue
					}
					err = ch.Publish("", pq.Name, false, false, amqp.Publishing{
						ContentType: "application/json",
						Body:        data,
					})
					if err != nil {
						t.logger.Error("pong publish error", zap.Error(err), zap.String("rabbit-queue", fmt.Sprintf("commandst.%s.pong", command)))
					}
				}
			}
		}
	}(msgs, t.channel, command)

	return nil
}
