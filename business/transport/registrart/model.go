package registrart

import "github.com/google/uuid"

type RegistrarSender interface {
	RegisterCommand(id uuid.UUID, command string) error
}

type RegistrarReceiver interface {
	ReceiveRegisterRequests() (<-chan RegistrarMessage, error)
	Close()
}

type RegistrarMessage struct {
	ID      uuid.UUID `json:"guild_id"`
	Command string    `json:"user_id"`
}
