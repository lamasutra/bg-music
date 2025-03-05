package ui

type UI interface {
	Debug(args ...any)
	Error(args ...any)
}

var Ready bool
var state UI

func CreateUI(ui string) UI {
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
