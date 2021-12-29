package service

import (
	"github.com/google/uuid"

	"github.com/FlameInTheDark/rebot/business/service/discovery"
)

type RegistrarHandler func(id uuid.UUID, command string)

type RegistrarWorker struct {
	registrar *discovery.ConsulDiscovery
	stop      chan struct{}
}

func NewRegistrarWorker(registrar *discovery.ConsulDiscovery) *RegistrarWorker {
	return &RegistrarWorker{
		registrar: registrar,
		stop:      make(chan struct{}),
	}
}

func (r *RegistrarWorker) AddRegistrarHandler(handler discovery.HandlerFunc) {
	r.registrar.AddHandler(handler)
}

func (r *RegistrarWorker) Run(name string) {
	go r.registrar.Start(name)
}
