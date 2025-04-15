package events

func (eb *eventBus) Listen(name string, key string, callback func(...any)) error {
	if !eb.Exists(name) {
		eb.registerEvent(name)
	}
	ev, err := eb.getEvent(name)
	if err != nil {
		return err
	}
	ev.listeners[key] = callback

	return nil
}

func (eb *eventBus) ListenAll(name string, callback globalEventListener) {
	eb.listeners[name] = callback
}

func Listen(name string, key string, callback func(...any)) error {
	return defaultBus.Listen(name, key, callback)
}

func ListenAll(name string, callback globalEventListener) {
	defaultBus.ListenAll(name, callback)
}
