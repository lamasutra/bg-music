package model

type BgPlayerConfig struct {
	Events  map[string]Event  `json:"events"`
	States  map[string]State  `json:"states"`
	Narrate map[string]Speech `json:"narrate"`
}

type ObjDistance struct {
	Combat uint32 `json:"combat"`
	Danger uint32 `json:"danger"`
	// Safe   uint32 `json:"safe"`
}
type Distances struct {
	Air    ObjDistance `json:"air"`
	Ground ObjDistance `json:"ground"`
}

type Theme struct {
	Title     string            `json:"title"`
	Events    map[string]Event  `json:"events"`
	States    map[string]State  `json:"states"`
	Narrate   map[string]Speech `json:"narrate"`
	Distances Distances         `json:"distances"`
	Extend    string            `json:"extend"`
}

func (d *Distances) merge(with Distances) *Distances {
	merged := Distances{}
	if with.Air.Combat > 0 {
		merged.Air.Combat = with.Air.Combat
	} else {
		merged.Air.Combat = d.Air.Combat
	}
	if with.Air.Danger > 0 {
		merged.Air.Danger = with.Air.Danger
	} else {
		merged.Air.Danger = d.Air.Danger
	}
	if with.Ground.Combat > 0 {
		merged.Ground.Combat = with.Ground.Combat
	} else {
		merged.Ground.Combat = d.Ground.Combat
	}
	if with.Ground.Danger > 0 {
		merged.Ground.Danger = with.Ground.Danger
	} else {
		merged.Ground.Danger = d.Ground.Danger
	}

	return &merged
}

func (t *Theme) mergeEvents(events map[string]Event) map[string]Event {
	var merged map[string]Event = make(map[string]Event)

	for stKey, event := range t.Events {
		merged[stKey] = event
	}

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

	for stKey, state := range t.States {
		merged[stKey] = state
	}

	for stKey, state := range states {
		destState, exists := merged[stKey]
		if exists {
			merged[stKey] = destState.merge(state)
		} else {
			merged[stKey] = state
		}
	}

	return merged
}

func (t *Theme) mergeNarrates(narrates map[string]Speech) map[string]Speech {
	var merged map[string]Speech = make(map[string]Speech)

	for nKey, speech := range t.Narrate {
		merged[nKey] = speech
	}

	for nKey, speech := range narrates {
		destNarrate, exists := merged[nKey]
		if exists {
			merged[nKey] = destNarrate.merge(speech)
		} else {
			merged[nKey] = speech
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
	merged := Theme{
		Events:    t.mergeEvents(with.Events),
		States:    t.mergeStates(with.States),
		Narrate:   t.mergeNarrates(with.Narrate),
		Distances: *t.Distances.merge(with.Distances),
	}

	if with.Title != "" {
		merged.Title = with.Title
	} else {
		merged.Title = t.Title
	}
	if with.Extend != "" {
		merged.Extend = with.Extend
	} else {
		merged.Extend = t.Extend
	}

	// fmt.Println("merged theme", utils.JsonPretty(merged))

	return &merged
}
