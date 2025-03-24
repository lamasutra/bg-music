package model

import (
	"github.com/gopxl/beep/v2"
)

type Player interface {
	Init()
	Play(stream beep.Streamer)
	PlayMusic(music *Music, c *Config) (beep.StreamSeekCloser, error)
	PlayMusicAtVolume(music *Music, c *Config, volume uint8) (beep.StreamSeekCloser, error)
	PlaySfx(sfx *Sfx, c *Config) (beep.StreamSeekCloser, error)
	Speak(sentence *[]Speech, c *Config)
	SetVolume(volume uint8)
	GetMusicEndedChan() chan (bool)
	GetCurrentMusic() *Music
	GetCurrentMusicProgress() float64
	Close()
}
