package registrart

import "github.com/google/uuid"

type RegistrarSender interface {
	RegisterCommand(id uuid.UUID, command string) error
}

type RegistrarReceiver interface {
	AddHandler(handler RegistrarHandler)
	Start() error
	Close()
}

type RegistrarMessage struct {
	ID      uuid.UUID `json:"guild_id"`
	Command string    `json:"user_id"`
}

type RegistrarHandler func(id uuid.UUID, command string)
