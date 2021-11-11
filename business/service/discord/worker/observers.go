package worker

import (
	"sync"

	"go.uber.org/zap"

	"github.com/google/uuid"
)

type ObserverType string

func (ot ObserverType) String() string {
	return string(ot)
}

type Observer interface {
	Notify(e *MessageEvent)
	Ping() error
	GetId() uuid.UUID
	Type() ObserverType
}

type MessageEvent struct {
	GuildID  string
	UserID   string
	Username string
	Message  string
}

type EventObservers struct {
	rw        sync.RWMutex
	observers []Observer
	logger    *zap.Logger
}

//Register add new observer
func (eo *EventObservers) Register(o Observer) {
	eo.rw.Lock()
	defer eo.rw.Unlock()
	if !eo.isObserverExists(o.GetId()) {
		eo.observers = append(eo.observers, o)
	}
}

func (eo *EventObservers) IsTypeExists(t ObserverType) bool {
	for _, e := range eo.observers {
		if e.Type() == t {
			return true
		}
	}
	return false
}

//Unregister remove observer
func (eo *EventObservers) Unregister(id uuid.UUID) {
	eo.rw.Lock()
	defer eo.rw.Unlock()
	for i, o := range eo.observers {
		if o.GetId() == id {
			eo.observers = append(eo.observers[:i], eo.observers[i+1:len(eo.observers)]...)
			return
		}
	}
}

//CheckHealth ping all observers and remove those who no longer respond
func (eo *EventObservers) CheckHealth() {
	wg := sync.WaitGroup{}
	for i, _ := range eo.observers {
		wg.Add(1)
		go func(o Observer) {
			if err := o.Ping(); err != nil {
				eo.logger.Warn("observer no longer responds", zap.String("observer-id", o.GetId().String()))
				eo.Unregister(o.GetId())
			}
			wg.Done()
		}(eo.observers[i])
	}
	wg.Wait()
}

//Notify observers about the event
func (eo *EventObservers) Notify(e *MessageEvent) {
	eo.rw.RLock()
	defer eo.rw.RUnlock()

	for _, o := range eo.observers {
		o.Notify(e)
	}
}

func (eo *EventObservers) isObserverExists(id uuid.UUID) bool {
	for _, o := range eo.observers {
		if o.GetId() == id {
			return true
		}
	}
	return false
}
