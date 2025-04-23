package input

type airGameMode struct {
}

func NewAirMode() *airGameMode {
	return &airGameMode{}
}

func (a *airGameMode) GetName() string {
	return "air"
}

func (a *airGameMode) ParseInput() {
}
