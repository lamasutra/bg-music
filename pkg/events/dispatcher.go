package events

func (eb *eventBus) dispatch(name string, values ...any) error {
	ev, err := eb.getEvent(name)
	if err != nil {
		return err
	}

	for _, callback := range ev.listeners {
		callback(values...)
	}

	for _, callback := range eb.listeners {
		callback(name, values...)
	}

	return nil
}
