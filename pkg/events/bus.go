package events

import (
	"fmt"
)

type eventListener func(...any)
type globalEventListener func(string, ...any)

type event struct {
	listeners map[string]eventListener
}

type eventBus struct {
	events    map[string]*event
	listeners map[string]globalEventListener
}

var defaultBus = &eventBus{
	events:    make(map[string]*event),
	listeners: make(map[string]globalEventListener),
}

func New() *eventBus {
	return &eventBus{
		events: make(map[string]*event),
	}
}

func (eb *eventBus) registerEvent(name string) error {
	if eb.Exists(name) {
		return fmt.Errorf("event `%s` already registered", name)
	}
	ev := event{
		listeners: make(map[string]eventListener),
	}

	eb.events[name] = &ev

	return nil
}

func (eb *eventBus) Exists(name string) bool {
	return eb.events[name] != nil
}

func (eb *eventBus) getEvent(name string) (*event, error) {
	if !eb.Exists(name) {
		return nil, fmt.Errorf("event `%s` is not registered", name)
	}

	event := eb.events[name]

	return event, nil
}

func Exists(name string) bool {
	return defaultBus.Exists(name)
}
