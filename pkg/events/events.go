package events

type eventListener func(...any)
type globalEventListener func(string, ...any)

type event struct {
	listeners map[string]eventListener
}

type eventBus struct {
	events    map[string]*event
	listeners map[string]globalEventListener
}
