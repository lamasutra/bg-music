package server

import "github.com/lamasutra/bg-music/model"

type Request struct {
	Action string `json:"action"`
}

type LoadData struct {
	Events map[string]model.Event `json:"events"`
	States map[string]model.State `json:"states"`
}

type LoadRequest struct {
	Request
	Data LoadData `json:"data"`
}
