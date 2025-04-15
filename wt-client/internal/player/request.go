package player

import "github.com/lamasutra/bg-music/wt-client/internal/model"

type Request struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type EventStates struct {
	Events map[string]model.Event `json:"events"`
	States map[string]model.State `json:"states"`
}
