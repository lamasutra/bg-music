package model

import (
	"github.com/gopxl/beep/v2"
)

type Player interface {
	Init()
	Play(stream beep.Streamer)
	SetPlaylist(playlist *[]Music)
	PlayMusic(music *Music, c *Config, allowSame bool)
	PlaySfx(sfx *Sfx, c *Config)
	Speak(sentence *[]Speech, c *Config)
	SetVolume(volume uint8)
	GetMusicEndedChan() chan (bool)
	GetCurrentMusic() *Music
	GetCurrentMusicProgress() float64
	SendControl(ctrl string)
	VolumeUp()
	VolumeDown()
	Next()
	Prev()
	Mute()
	Pause()
	Close()
}
