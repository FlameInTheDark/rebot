package observers

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/service/discord/worker"
	"github.com/FlameInTheDark/rebot/business/transport/commandst"
)

type RabbitCommandObserver struct {
	id           uuid.UUID
	event        string
	observerType worker.ObserverType
	cmdSender    commandst.CommandsSender

	logger *zap.Logger
}

func NewRabbitCommandObserver(id uuid.UUID, event string, cmdSender commandst.CommandsSender, ot worker.ObserverType, logger *zap.Logger) *RabbitCommandObserver {
	return &RabbitCommandObserver{
		id:           id,
		event:        event,
		observerType: ot,
		cmdSender:    cmdSender,
		logger:       logger,
	}
}

func (r *RabbitCommandObserver) Notify(e *worker.MessageEvent) {
	msg := commandst.CommandMessage{
		GuildID:  e.GuildID,
		UserID:   e.UserID,
		Username: e.Username,
		Message:  e.Message,
	}
	err := r.cmdSender.SendCommand(msg, r.event)
	if err != nil {
		r.logger.Error(
			"cannot send event message",
			zap.Error(err),
			zap.String("event-observer-id", r.id.String()),
			zap.String("event-observer-type", r.observerType.String()),
		)
	}
}

func (r *RabbitCommandObserver) Ping() error {
	return r.cmdSender.Ping(r.event)
}

func (r *RabbitCommandObserver) GetId() uuid.UUID {
	return r.id
}

func (r *RabbitCommandObserver) Type() worker.ObserverType {
	return r.observerType
}
