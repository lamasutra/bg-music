package player

import (
	"github.com/gopxl/beep/v2"
	"github.com/lamasutra/bg-music/model"
)

var player model.Player

func CreatePlayer(playerType string) *model.Player {

	switch playerType {
	case "beep":
		player = &(beepState{
			format: beep.Format{
				SampleRate:  44100,
				NumChannels: 2,
				Precision:   2,
			},
			sequencers: sequencers{
				music:    NewBeepSequencer(8),
				sfx:      NewBeepSequencer(8),
				narrator: NewBeepSequencer(32),
			},
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
