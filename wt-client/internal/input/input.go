package input

import (
	"github.com/lamasutra/bg-music/wt-client/internal/model"
)

type GameMode interface {
	GetName() string
	ParseInput()
}

type DataParser struct {
	modes       map[string]GameMode
	inputData   *model.WtData
	currentMode GameMode
}

var currentVehicle *model.Vehicle
var currentTheme *model.Theme

var _inputMapBool_ model.WtInputMapBool = make(model.WtInputMapBool, 9)
var inputMapBool *model.WtInputMapBool = &_inputMapBool_
var input = &model.WtInput{
	NearestEnemyAir:    -1,
	NearestEnemyGround: -1,
}

var shouldStayDead bool
