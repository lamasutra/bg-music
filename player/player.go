package player

import (
	"github.com/lamasutra/bg-music/model"
)

var player model.Player

func CreatePlayer(playerType string) *model.Player {

	switch playerType {
	case "beep":
		player = &(beepState{
			// volumePercent: volume,
			// musicEnded:    musicEndedChannel,
			// stopWatchEnd:  make(chan bool, 1),
		})
	}

	if player == nil {
		return nil
	}
	player.Init()

	return &player
}

func GetCurrentPlayer() *model.Player {
	return &player
}
