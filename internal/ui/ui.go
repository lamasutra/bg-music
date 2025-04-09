package ui

import (
	"embed"

	"github.com/lamasutra/bg-music/model"
)

var state model.UI

func CreateUI(ui string, assets *embed.FS, icon []byte, onStartup func()) {
	switch ui {
	case "tui":
		state = NewTui()
	case "cli":
		state = NewCli()
	default:
		state = NewGui(assets, icon)
	}

	state.Run(onStartup)
}

func Debug(args ...any) {
	state.Debug(args...)
}

func Error(args ...any) {
	state.Error(args...)
}

func GetState() *model.UI {
	return &state
}
