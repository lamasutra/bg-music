package input

type tanksGameMode struct {
}

func NewTanksMode() *tanksGameMode {
	return &tanksGameMode{}
}

func (t *tanksGameMode) GetName() string {
	return "tanks"
}

func (t *tanksGameMode) ParseInput() {
}
