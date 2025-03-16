package ui

import "github.com/lamasutra/bg-music/model"

var state model.UI

func CreateUI(ui string) {
	switch ui {
	case "tui":
		state = NewTui()
	default:
		state = NewCli()
	}
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

func GetState() *model.UI {
	return &state
}
