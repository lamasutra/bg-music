package player

import "github.com/lamasutra/bg-music/wt-client/clientConfig"

type BgPlayer interface {
	Init(*clientConfig.Config)
	SendEventStates(*clientConfig.EventStates) error
	Close()
}

func CreatePlayer(playerType string, c *clientConfig.Config) BgPlayer {
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
