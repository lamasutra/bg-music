package model

type EventStates struct {
	Events map[string]Event `json:"events"`
	States map[string]State `json:"states"`
}

type Theme struct {
	Title  string
	Events map[string]Event `json:"events"`
	States map[string]State `json:"states"`
}

func (t *Theme) mergeEvents(events map[string]Event) map[string]Event {
	var merged map[string]Event

	for evKey, event := range t.Events {
		destEvent, exists := events[evKey]
		if exists {
			merged[evKey] = destEvent.merge(event)
		} else {
			merged[evKey] = event
		}
	}

	return merged
}

func (t *Theme) mergeStates(states map[string]State) map[string]State {
	var merged map[string]State

	for stKey, state := range t.States {
		destState, exists := states[stKey]
		if exists {
			merged[stKey] = destState.merge(state)
		} else {
			merged[stKey] = state
		}
	}

	return merged
}

func (t *Theme) Merge(with Theme) *Theme {
	merged := &Theme{
		Title:  with.Title,
		Events: t.mergeEvents(with.Events),
		States: t.mergeStates(with.States),
	}

	return merged
}
