package events

func (eb *eventBus) Trigger(name string, values ...any) {
	if !eb.Exists(name) {
		eb.registerEvent(name)
	}
	eb.dispatch(name, values...)
}

func Trigger(event string, values ...any) {
	defaultBus.Trigger(event, values...)
}
