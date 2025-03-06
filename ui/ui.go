package ui

import "github.com/lamasutra/bg-music/model"

var Ready bool
var state model.UI

func CreateUI(ui string) model.UI {
	switch ui {
	case "tui":
		state = NewTui()
	default:
		state = NewCli()
	}

	Ready = true

	return state
}

func Debug(args ...any) {
	state.Debug(args...)
}

func Error(args ...any) {
	state.Error(args...)
}

func SetCurrentMusicTitle(title string) {
	state.SetCurrentMusicTitle(title)
}

func SetCurrentMusicProgress(value float64) {
	// Debug("progress", value)
	state.SetCurrentMusicProgress(value)
}

func SetCurrentVolume(value float64) {
	// Debug("progress", value)
	state.SetCurrentVolume(value)
}
