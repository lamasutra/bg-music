package player

import "github.com/lamasutra/bg-music/wt-client/model"

type BgPlayer interface {
	Init(*model.Config)
	SendEventStates(*model.BgPlayerConfig) error
	ChangeMusic() error
	TriggerEvent(string) error
	SendState(string) error
	Speak(string) error
	Close()
}

func CreatePlayer(playerType string, c *model.Config) BgPlayer {
	var player BgPlayer
	switch playerType {
	case "pipe":
		player = &PipePlayer{}
	case "http":
		player = &HttPlayer{}
	default:
		panic("unknown player")
	}

	player.Init(c)

	return player
}
