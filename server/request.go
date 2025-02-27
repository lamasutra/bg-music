package server

import "github.com/lamasutra/bg-music/model"

type Request struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type LoadConfigData struct {
	Events map[string]model.Event `json:"events"`
	States map[string]model.State `json:"states"`
}
