package discovert

import "github.com/google/uuid"

type DiscoverSender interface {
	DiscoverCommand() error
}

type DiscoverReceiver interface {
	ReceiveDiscoverRequests() (<-chan DiscoverMessage, error)
	DiscoverResponse() error
	Close()
}

type DiscoverPing string

const PingMessage = DiscoverPing("ping")

type DiscoverMessage struct {
	ID      uuid.UUID `json:"guild_id"`
	Command string    `json:"user_id"`
}
