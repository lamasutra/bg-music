package ui

type UI interface {
}

func CreateUI() UI {
	return NewTui()
}
