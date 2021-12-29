package worker

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/transport/commandst"
)

type Observer interface {
	Notify(e *MessageEvent)
}

type MessageEvent struct {
	GuildID   string
	UserID    string
	ChannelID string
	Username  string
	Message   string
}

type RabbitCommandObserver struct {
	id        uuid.UUID
	event     string
	cmdSender commandst.CommandsSender

	logger *zap.Logger
}

func NewRabbitCommandObserver(event string, cmdSender commandst.CommandsSender, logger *zap.Logger) *RabbitCommandObserver {
	return &RabbitCommandObserver{
		event:     event,
		cmdSender: cmdSender,
		logger:    logger,
	}
}

func (r *RabbitCommandObserver) Notify(e *MessageEvent) {
	r.logger.Debug("Sending event message", zap.Reflect("rabbit-message", e))
	msg := commandst.CommandMessage{
		GuildID:   e.GuildID,
		ChannelID: e.ChannelID,
		UserID:    e.UserID,
		Username:  e.Username,
		Message:   e.Message,
	}
	err := r.cmdSender.SendCommand(msg, r.event)
	if err != nil {
		r.logger.Error(
			"cannot send event message",
			zap.Error(err),
			zap.String("event-observer-id", r.id.String()),
		)
	}
}

func (r *RabbitCommandObserver) GetId() uuid.UUID {
	return r.id
}
