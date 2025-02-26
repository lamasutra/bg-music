package model

type State struct {
	Event
	States []string `json:"states"`
}
