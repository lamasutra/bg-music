package ui

type UI interface {
}

var state UI

func CreateUI(ui string) {
	switch ui {
	case "cli":
		state = NewCli()
	case "tui":
		state = NewTui()
	default:
		// state = NewGui()
		state = NewCli()
	}
}
