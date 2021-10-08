package worker

import (
	"github.com/google/uuid"
	"sync"
)

type Observer interface {
	Use(e *Event)
	Ping() error
	GetId() uuid.UUID
}

type Event struct {
	GuildID string
	UserID  string
	Message string
}

type DiscordEvents struct {
	rw        sync.RWMutex
	observers []Observer
}

//Register add new observer
func (de *DiscordEvents) Register(o Observer) {
	de.rw.Lock()
	defer de.rw.Unlock()
	if !de.isObserverExists(o.GetId()) {
		de.observers = append(de.observers, o)
	}
}

func (de *DiscordEvents) Unregister(id uuid.UUID) {
	de.rw.Lock()
	defer de.rw.Unlock()
	for i,o := range de.observers {
		if o.GetId() == id {
			de.observers = append(de.observers[:i], de.observers[i+1:len(de.observers)]...)
			return
		}
	}
}

func (de *DiscordEvents) Use(e *Event) {
	de.rw.RLock()
	defer de.rw.RUnlock()

	for _,o := range de.observers {
		o.Use(e)
	}
}

func (de *DiscordEvents) isObserverExists(id uuid.UUID) bool {
	for _,o := range de.observers {
		if o.GetId() == id {
			return true
		}
	}
	return false
}