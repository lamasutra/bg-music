package player

import (
	"github.com/gopxl/beep/v2"
	"github.com/lamasutra/bg-music/config"
	"github.com/lamasutra/bg-music/model"
)

type Player interface {
	Init()
	PlayMusic(music *model.Music, c *config.Config) (beep.StreamSeekCloser, error)
	PlaySfx(sfx *model.Sfx, c *config.Config) (beep.StreamSeekCloser, error)
	SetVolume(volume uint8)
	Close()
}

func CreatePlayer(playerType string, volume uint8) Player {
	var player Player
	switch playerType {
	case "beep":
		player = &(BeepPlayer{
			volumePercent: volume,
		})
	}

	if player == nil {
		return nil
	}
	player.Init()

	return player
}
