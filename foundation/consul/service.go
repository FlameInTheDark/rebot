package consul

import "github.com/google/uuid"

type Service struct {
	ID   uuid.UUID         `json:"id"`
	Name string            `json:"name"`
	Tags []string          `json:"tags"`
	Meta map[string]string `json:"meta"`
}
