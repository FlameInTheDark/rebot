package worker

import (
	"go.uber.org/zap"
	"sync"
)

// EventRegistrar contains events and subscribes observers on events
type EventRegistrar struct {
	rw     sync.RWMutex
	events map[string]*EventObservers
	logger *zap.Logger
}

//NewEventRegistrar creates a new EventRegistrar
func NewEventRegistrar() *EventRegistrar {
	return &EventRegistrar{
		events: make(map[string]*EventObservers),
	}
}

//Register create new event observer and register observer
func (er *EventRegistrar) Register(event string, observer Observer) {
	er.rw.Lock()
	defer er.rw.Unlock()
	if _, ok := er.events[event]; !ok {
		er.events[event] = &EventObservers{logger: er.logger.With(zap.String("event", event))}
	}
	er.events[event].Register(observer)
}

// Notify triggers the event observer
func (er *EventRegistrar) Notify(event string, message MessageEvent) {
	er.rw.RLock()
	er.events[event].Notify(&message)
	er.rw.RUnlock()
}

// IsObserverTypeExist returns true if observer with given type already exists
func (er *EventRegistrar) IsObserverTypeExist(t ObserverType) bool {
	er.rw.RLock()
	defer er.rw.RUnlock()
	for _, e := range er.events {
		if e.IsTypeExists(t) {
			return true
		}
	}
	return false
}

//CheckObserversHealth check observers health
func (er *EventRegistrar) CheckObserversHealth() {
	er.rw.RLock()
	defer er.rw.RUnlock()
	for _, e := range er.events {
		e.CheckHealth()
	}
}
