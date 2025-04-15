package ui

import (
	"embed"
)

type UI interface {
	Run(func())
}

var ui UI

func CreateNew(uiType string, assets *embed.FS, icon []byte, onStartup func()) {
	switch uiType {
	case "tui":
		ui = NewTui()
	case "cli":
		ui = NewCli()
	default:
		ui = NewGui(assets, icon)
	}

	ui.Run(onStartup)
}
