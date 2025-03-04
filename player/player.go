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
	GetMusicEndedChan() *chan (bool)
	Close()
}

func CreatePlayer(playerType string, volume uint8, musicEndedChannel *chan bool) Player {
	var player Player
	switch playerType {
	case "beep":
		player = &(BeepPlayer{
			volumePercent: volume,
			musicEnded:    musicEndedChannel,
		})
	}

	if player == nil {
		return nil
	}
	player.Init()

	return player
}
