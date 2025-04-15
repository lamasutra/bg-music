package input

import (
	"time"

	"github.com/lamasutra/bg-music/wt-client/internal/client"
	"github.com/lamasutra/bg-music/wt-client/internal/model"
	"github.com/lamasutra/bg-music/wt-client/internal/types"
)

const sleepTime = time.Millisecond * 500
const sleepOffline = time.Millisecond * 1000

var currentVehicle *model.Vehicle
var currentTheme *model.Theme

var inputData = &types.WtData{
	State:      &client.State{},
	MapInfo:    &client.MapInfo{},
	MapObj:     &client.MapObj{},
	Indicators: &client.Indicators{},
	HudMsg:     &client.HudMsg{},
}

var state struct {
	lastEvt uint64
	lastDmg uint64
}

var _inputMapBool_ types.WtInputMapBool = make(types.WtInputMapBool, 9)
var inputMapBool *types.WtInputMapBool = &_inputMapBool_
var input = &types.WtInput{
	NearestEnemyAir:    -1,
	NearestEnemyGround: -1,
}

var shouldStayDead bool
