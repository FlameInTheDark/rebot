package worker

import (
	"go.uber.org/zap"
	"sync"
)

// EventRegistrar contains events and subscribes observers on events
type EventRegistrar struct {
	rw     sync.RWMutex
	events map[string]Observer
	logger *zap.Logger
}

//NewEventRegistrar creates a new EventRegistrar
func NewEventRegistrar(logger *zap.Logger) *EventRegistrar {
	return &EventRegistrar{
		logger: logger,
		events: make(map[string]Observer),
	}
}

//Register create new event observer and register observer
func (er *EventRegistrar) Register(event string, observer Observer) {
	er.rw.Lock()
	defer er.rw.Unlock()
	if _, ok := er.events[event]; !ok {
		er.events[event] = observer
	}
}

// Notify triggers the event observer
func (er *EventRegistrar) Notify(event string, message MessageEvent) {
	er.rw.RLock()
	defer er.rw.RUnlock()
	if _, ok := er.events[event]; !ok {
		return
	}
	er.events[event].Notify(&message)
}
