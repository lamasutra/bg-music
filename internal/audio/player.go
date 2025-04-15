package audio

import (
	"github.com/gopxl/beep/v2"
	"github.com/lamasutra/bg-music/pkg/model"
)

type Player interface {
	Init()
	Play(stream beep.Streamer)
	SetPlaylist(playlist *[]model.Music)
	PlayMusic(music *model.Music, c *model.Config, allowSame bool)
	PlaySfx(sfx *model.Sfx, c *model.Config)
	Speak(sentence *[]model.Speech, c *model.Config)
	SetVolume(volume uint8)
	GetMusicEndedChan() chan (bool)
	GetCurrentMusic() *model.Music
	GetCurrentMetadata() *model.MusicMetadata
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

var player Player

func CreatePlayer(playerType string) Player {

	switch playerType {
	case "beep":
		player = CreateBeepPlayer()
	}

	if player == nil {
		return nil
	}
	player.Init()

	return player
}

func GetCurrentPlayer() Player {
	return player
}
