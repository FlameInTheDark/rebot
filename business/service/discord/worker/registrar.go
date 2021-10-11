package worker

import "go.uber.org/zap"

type EventRegistrar struct {
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
	if _, ok := er.events[event]; !ok {
		er.events[event] = &EventObservers{logger: er.logger.With(zap.String("event", event))}
	}
	er.events[event].Register(observer)
}

//CheckObserversHealth check observers health
func (er *EventRegistrar) CheckObserversHealth() {
	for _,e := range er.events {
		e.CheckHealth()
	}
}
