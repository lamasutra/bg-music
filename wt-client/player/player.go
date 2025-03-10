package player

import "github.com/lamasutra/bg-music/wt-client/model"

type BgPlayer interface {
	Init(*model.Config)
	SendEventStates(*model.EventStates) error
	TriggerEvent(event string) error
	SendState(event string) error
	Close()
}

func CreatePlayer(playerType string, c *model.Config) BgPlayer {
	var player BgPlayer
	switch playerType {
	case "pipe":
		player = &PipePlayer{}
	default:
		panic("unknown player")
	}

	player.Init(c)

	return player
}
