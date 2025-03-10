package input

import (
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/player"
	"github.com/lamasutra/bg-music/wt-client/stateMachine"
	"github.com/lamasutra/bg-music/wt-client/types"
)

const sleepTime = time.Millisecond * 500
const sleepOffline = time.Millisecond * 1000

var state_ts int64

var inputData = &types.WtData{
	State:      &client.State{},
	MapInfo:    &client.MapInfo{},
	MapObj:     &client.MapObj{},
	Indicators: &client.Indicators{},
	HudMsg:     &client.HudMsg{},
}

var _inputMapBool_ types.WtInputMapBool = make(types.WtInputMapBool, 9)
var inputMapBool *types.WtInputMapBool = &_inputMapBool_
var input = &types.WtInput{}
var lastEvt, lastDmg uint64

// @todo paralel load
func loadData(host string) {
	// fmt.Println("Loading data from", host)
	err := inputData.State.Load(host)
	if err != nil {
		// fmt.Println("state error: ", err)
		time.Sleep(sleepOffline)
		err = inputData.State.Load(host)
		if err != nil {
			input.GameRunning = false
			(*inputMapBool)["GameRunning"] = input.GameRunning
			return
		} else {
			input.GameRunning = true
		}
	} else {
		input.GameRunning = true
		// @todo fetch lastMsg
	}

	(*inputMapBool)["GameRunning"] = input.GameRunning

	err = inputData.Indicators.Load(host)
	if err != nil {
		fmt.Println("indicators error: ", err)
	}
	err = inputData.MapInfo.Load(host)
	if err != nil {
		fmt.Println("mapInfo error: ", err)
	}
	// load other data
	if inputData.MapInfo.Valid {
		// load map identity
		if inputData.Identity == 0 {
			inputData.Identity, err = client.MapIdentity(host)
			fmt.Println("map identity error: ", err)
		}

		err = inputData.MapObj.Load(host)
		if err != nil {
			fmt.Println("mapObj error: ", err)
			// } else {
			// fmt.Println(mapObj)
		}
		err = inputData.HudMsg.Load(host, lastEvt, lastDmg)
		if err != nil {
			fmt.Println("hudMsg error: ", err)
		}
	} else {
		if inputData.Identity != 0 {
			inputData.Identity = 0
		}
	}
}

func parseInput(conf *model.Config) {
	// fmt.Println("Data:", inputData)

	input.MapLoaded = inputData.MapInfo.Valid

	if input.MapLoaded {
		player := inputData.MapObj.GetPlayerEntity()
		enemyAircrafts := inputData.MapObj.GetEnemyAircrafts(&conf.Colors.Enemy.Air)
		enemyGroundUnits := inputData.MapObj.GetEnemyGroundUnits(&conf.Colors.Enemy.Ground)
		isShotDown := inputData.HudMsg.IsShotDown(conf.Nickname)
		hasCrashed := inputData.HudMsg.HasCrashed(conf.Nickname)
		// isMissionEnded := inputData.hudMsg.IsMissionEnded()
		input.MissionEnded = false
		if !input.MissionStarted {
			// set lastDamage to the latest on mission start
			lastDamage := inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				lastDmg = uint64(lastDamage.ID)
			}
			input.MissionStarted = player != nil // @todo what else ?
		}
		input.PlayerLanded = input.MissionStarted && player == nil && !isShotDown && !hasCrashed
		input.PlayerDead = input.MissionStarted && player == nil && (isShotDown || hasCrashed)
		// fmt.Println("dead conds", input.PlayerDead, isShotDown, hasCrashed)
		input.PlayerType = inputData.Indicators.Army
		input.PlayerVehicle = inputData.Indicators.Type
		// @todo last known location
		if player != nil {
			nearestAir := getNearestEntity(player, enemyAircrafts, inputData.MapObj, inputData.MapInfo)
			nearestGround := getNearestEntity(player, enemyGroundUnits, inputData.MapObj, inputData.MapInfo)

			// @todo configurable in theme
			if nearestAir != nil {
				input.EnemyAirNear = inputData.MapObj.GetDistance(player, nearestAir, inputData.MapInfo) < 10000
				input.EnemyAirClose = inputData.MapObj.GetDistance(player, nearestAir, inputData.MapInfo) < 5000
			}

			// @todo configurable in theme
			if nearestGround != nil {
				input.EnemyGroundNear = inputData.MapObj.GetDistance(player, nearestGround, inputData.MapInfo) < 20000
				input.EnemyGroundClose = inputData.MapObj.GetDistance(player, nearestGround, inputData.MapInfo) < 1000
			}

			// @todo enemy base
		} else {
			input.EnemyAirNear = false
			input.EnemyAirClose = false
		}

		// only set lastDamage if player is dead and resurrected
		// if input.PlayerDead && player != nil {
		lastDamage := inputData.HudMsg.GetLastDmg()
		if lastDamage != nil {
			lastDmg = uint64(lastDamage.ID)
		}
		// }
	} else {
		input.MapLoaded = false

		if input.PlayerType != "" {
			input.MissionEnded = true
		}

		input.MissionStarted = false
		input.PlayerDead = false
		input.PlayerLanded = false
		input.EnemyAirClose = false
		input.EnemyAirNear = false
		input.EnemyGroundClose = false
		input.EnemyGroundNear = false
		input.PlayerType = ""
		input.PlayerVehicle = ""
	}

	(*inputMapBool)["MapLoaded"] = input.MapLoaded
	(*inputMapBool)["MapLoaded"] = input.MapLoaded
	(*inputMapBool)["MissionEnded"] = input.MissionEnded
	(*inputMapBool)["MissionStarted"] = input.MissionStarted
	(*inputMapBool)["PlayerDead"] = input.PlayerDead
	(*inputMapBool)["PlayerLanded"] = input.PlayerLanded
	(*inputMapBool)["EnemyAirClose"] = input.EnemyAirClose
	(*inputMapBool)["EnemyAirNear"] = input.EnemyAirNear
	(*inputMapBool)["EnemyGroundClose"] = input.EnemyGroundClose
	(*inputMapBool)["EnemyGroundNear"] = input.EnemyGroundNear

	// buf, _ := json.MarshalIndent(input, "", "  ")
	// fmt.Println("Input", string(buf))

	// fmt.Println("last", lastDmg, lastEvt)

	// var state string
	// if input.GameRunning {
	// 	state = "running"
	// } else {
	// 	state = "offline"
	// }
	// fmt.Println("Game:", state)
}

func getNearestEntity(player *client.Entity, entities *[]client.Entity, mapObj *client.MapObj, mapInfo *client.MapInfo) *client.Entity {
	var nearest *client.Entity
	var distance, lastDistance float64
	for _, entity := range *entities {
		distance = mapObj.GetDistance(player, &entity, mapInfo)
		if lastDistance == 0.0 || distance < lastDistance {
			lastDistance = distance
			nearest = &entity
		}
	}

	return nearest
}

func LoadLoop(host string, conf *model.Config, stMachine *stateMachine.StateMachine, player player.BgPlayer) {
	// fmt.Println(inputData)

	currentState := stMachine.GetCurrentState()
	var newState string
	var current_ts int64
	var cooldown_s int64

	// @todo find recent state
	fmt.Print("sending default state ", currentState, " ... ")

	err := player.SendState(currentState)
	if err != nil {
		fmt.Println("failed")
	} else {
		fmt.Println("ok")
	}

	for {
		loadData(host)
		parseInput(conf)
		current_ts = time.Now().Unix()
		// fast forward state
		newState = ""
		for {
			state, err := stMachine.GetNextState(inputMapBool)
			if err != nil {
				// fmt.Println("getNextState failed", err)
				break
			}
			// state change cooldown, do not set state in colldown period
			if current_ts < (state_ts + cooldown_s) {
				break
			}
			// fmt.Println("st", state, stMachine.GetCurrentState())
			if state != "" {
				newState = state
				stMachine.SetState(state)
				state_ts = time.Now().Unix()
				// fmt.Println("state", state)
			} else {
				if newState != "" {
					// fmt.Println("new state:", newState)
					player.SendState(newState)
					cooldown_s = conf.
					// } else {
					// fmt.Println("state not changed")
					// }
				}
				break
			}
			time.Sleep(time.Millisecond * 100)
		}

		// fmt.Println(inputMapBool)
		time.Sleep(sleepTime)
	}
}

func Close() {
	// close(chRead)
	// close(chInput)
}
