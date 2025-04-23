package ui

import "github.com/lamasutra/bg-music/internal/app"

type UI interface {
	Run(func(a *app.AppState))
}

var ui UI

func CreateNew(uiType string, a *app.AppState, onStartup func(a *app.AppState)) {
	switch uiType {
	case "tui":
		ui = NewTui(a)
	case "cli":
		ui = NewCli(a)
	default:
		ui = NewGui(a)
	}

	ui.Run(onStartup)
}
