package player

import "github.com/lamasutra/bg-music/wt-client/clientConfig"

type Request struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type EventStates struct {
	Events map[string]clientConfig.Event `json:"events"`
	States map[string]clientConfig.State `json:"states"`
}
