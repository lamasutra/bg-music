package model

import (
	"github.com/gopxl/beep/v2"
)

type Player interface {
	Init()
	PlayMusic(music *Music, c *Config) (beep.StreamSeekCloser, error)
	PlaySfx(sfx *Sfx, c *Config) (beep.StreamSeekCloser, error)
	SetVolume(volume uint8)
	GetMusicEndedChan() chan (bool)
	GetCurrentMusic() *Music
	GetCurrentMusicProgress() float64
	Close()
}
