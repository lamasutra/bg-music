package player

import (
	"github.com/bg-music/config"
	"github.com/bg-music/model"
	"github.com/gopxl/beep/v2"
)

type Player interface {
	PlayMusic(music *model.Music, c *config.Config) (beep.StreamSeekCloser, error)
	PlaySfx(sfx *model.Sfx, c *config.Config) (beep.StreamSeekCloser, error)
	SetVolume(volume uint8)
	Close()
}

func CreatePlayer(playerType string, volume uint8) Player {
	switch playerType {
	case "local":
		return &(LocalPlayer{
			volumePercent: volume,
		})
	}

	return nil
}
