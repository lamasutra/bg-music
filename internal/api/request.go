package api

import "github.com/lamasutra/bg-music/pkg/model"

type Request struct {
	Action string `json:"action"`
}

type LoadData struct {
	Events  map[string]model.Event  `json:"events"`
	States  map[string]model.State  `json:"states"`
	Narrate map[string]model.Speech `json:"narrate"`
}

type LoadRequest struct {
	Request
	Data LoadData `json:"data"`
}

type StateRequest struct {
	State string `json:"state"`
}

type EventRequest struct {
	Event string `json:"event"`
}
