package ui

import "github.com/lamasutra/bg-music/wt-client/types"

type UI interface {
	Debug(...any)
	Error(...any)
	Input(*types.WtInput)
}

var state UI

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

func Input(in *types.WtInput) {
	state.Input(in)
}
