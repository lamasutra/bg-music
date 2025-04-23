package input

type defaultGameMode struct {
}

func NewDefaultMode() *defaultGameMode {
	return &defaultGameMode{}
}

func (d *defaultGameMode) GetName() string {
	return "default"
}

func (d *defaultGameMode) ParseInput() {
	// just dummy parser
}
