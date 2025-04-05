package model

import "io"

type UI interface {
	io.Writer
	Debug(...any)
	Error(...any)
	SetCurrentMusicTitle(string)
	SetCurrentMusicProgress(float64)
	SetCurrentVolume(float64)
}
