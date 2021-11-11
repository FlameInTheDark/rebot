package discovert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var _ DiscoverSender = (*RabbitDiscoverTransport)(nil)
var _ DiscoverReceiver = (*RabbitDiscoverTransport)(nil)

type RabbitDiscoverTransport struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	logger *zap.Logger

	close chan struct{}
}

func NewRabbitDiscoverTransport(conn *amqp.Connection, logger *zap.Logger) (*RabbitDiscoverTransport, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create channel")
	}

	return &RabbitDiscoverTransport{
		conn:    conn,
		channel: ch,
		logger:  logger,
		close:   make(chan struct{}),
	}, nil
}

func (t *RabbitDiscoverTransport) DiscoverCommand() error {
	discq, err := t.getQueue("service.discover.ping")
	if err != nil {
		return err
	}

	var ping = PingMessage
	data, err := json.Marshal(&ping)
	if err != nil {
		return err
	}

	err = t.channel.Publish("", discq.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *RabbitDiscoverTransport) Close() {
	close(t.close)
}

func (t *RabbitDiscoverTransport) getQueue(name string) (amqp.Queue, error) {
	return t.channel.QueueDeclare(name, true, true, false, false, nil)
}

func (t *RabbitDiscoverTransport) SendCommand(cmd CommandMessage, command string) error {
	data, err := json.Marshal(&cmd)
	if err != nil {
		return err
	}
	q, err := t.getQueue(fmt.Sprintf("commands.%s", command))
	if err != nil {
		return err
	}
	return t.channel.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func (t *RabbitDiscoverTransport) ReceiveCommands(command string) (<-chan CommandMessage, error) {
	var ch = make(chan CommandMessage)

	q, err := t.getQueue(fmt.Sprintf("commands.%s", command))
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

func (t *RabbitDiscoverTransport) ReceivePings(command string) error {
	q, err := t.getQueue(fmt.Sprintf("commands.%s.ping", command))
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
						fmt.Sprintf("commands.%s.pong", command),
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
