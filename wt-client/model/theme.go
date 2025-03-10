package model

import (
	"fmt"

	"github.com/lamasutra/bg-music/wt-client/utils"
)

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
	var merged map[string]Event = make(map[string]Event)

	for evKey, event := range events {
		destEvent, exists := t.Events[evKey]
		if exists {
			merged[evKey] = destEvent.merge(event)
		} else {
			merged[evKey] = event
		}
	}

	return merged
}

func (t *Theme) mergeStates(states map[string]State) map[string]State {
	var merged map[string]State = make(map[string]State)

	for stKey, state := range states {
		destState, exists := t.States[stKey]
		if exists {
			merged[stKey] = destState.merge(state)
		} else {
			merged[stKey] = state
		}
	}

	return merged
}

func (t *Theme) forceStateVolume(volume int) map[string]State {
	states := make(map[string]State, len(t.States))

	for idx, state := range t.States {
		state.Volume = volume
		states[idx] = state
	}

	return states
}

func (t *Theme) merge(with Theme) *Theme {
	merged := &Theme{
		Title:  with.Title,
		Events: t.mergeEvents(with.Events),
		States: t.mergeStates(with.States),
	}

	fmt.Println("merged theme", utils.JsonPretty(merged))

	return merged
}
