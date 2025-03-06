package model

type UI interface {
	Debug(...any)
	Error(...any)
	SetCurrentMusicTitle(string)
	SetCurrentMusicProgress(float64)
	SetCurrentVolume(float64)
}
