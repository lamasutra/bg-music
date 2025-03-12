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
var currentVehicle *model.Vehicle
var currentTheme *model.Theme

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
	var objDistance float64

	if input.MapLoaded {
		player := inputData.MapObj.GetPlayerEntity()
		enemyAircrafts := inputData.MapObj.GetEnemyAircrafts(&conf.Colors.Enemy.Air)
		enemyGroundUnits := inputData.MapObj.GetEnemyGroundUnits(&conf.Colors.Enemy.Ground)
		isShotDown := inputData.HudMsg.IsShotDown(conf.Nickname)
		hasCrashed := inputData.HudMsg.HasCrashed(conf.Nickname)
		playerScored := inputData.HudMsg.HadKill(conf.Nickname)
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

		// vehicle changed
		if input.PlayerVehicle != inputData.Indicators.Type {
			input.PlayerType = inputData.Indicators.Army
			input.PlayerVehicle = inputData.Indicators.Type
			currentVehicle = conf.GetVehicleForPlayerTypeAndVehicleTitle(input.PlayerType, input.PlayerVehicle)
			currentTheme = conf.GetThemeForVehicle(currentVehicle)
		}

		// @todo last known location
		if player != nil {
			nearestAir := getNearestEntity(player, enemyAircrafts, inputData.MapObj, inputData.MapInfo)
			nearestGround := getNearestEntity(player, enemyGroundUnits, inputData.MapObj, inputData.MapInfo)

			// @todo configurable in theme
			if nearestAir != nil {
				objDistance = inputData.MapObj.GetDistance(player, nearestAir, inputData.MapInfo)
				// 10000
				input.EnemyAirNear = objDistance < float64(currentTheme.Distances.Air.Danger)
				// 5000
				input.EnemyAirClose = objDistance < float64(currentTheme.Distances.Air.Combat)
				// fmt.Println("air dist", objDistance, float64(currentTheme.Distances.Air.Danger), float64(currentTheme.Distances.Air.Combat), input.EnemyAirNear, input.EnemyAirClose)
			} else {
				input.EnemyAirNear = false
				input.EnemyAirClose = false
			}

			// @todo configurable in theme
			if nearestGround != nil {
				inputData.MapObj.GetDistance(player, nearestGround, inputData.MapInfo)
				// 20000
				input.EnemyGroundNear = objDistance < float64(currentTheme.Distances.Ground.Danger)
				// 1000
				input.EnemyGroundClose = objDistance < float64(currentTheme.Distances.Ground.Combat)
			} else {
				// fmt.Println("ground dist", objDistance, float64(currentTheme.Distances.Ground.Danger), float64(currentTheme.Distances.Ground.Combat))
				input.EnemyGroundNear = false
				input.EnemyGroundClose = false
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
		currentVehicle = nil
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

func getCurrentVehicle() *model.Vehicle {
	return currentVehicle
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

func LoadLoop(host string, conf *model.Config, stMachine *stateMachine.StateMachine, bgPlayer player.BgPlayer) {
	fmt.Println(inputData)

	currentState := stMachine.GetCurrentState()
	var newState, state string
	var current_ts int64
	var cooldown_s int64
	var currentVehicle string
	var vehicleConf *model.Vehicle
	var vehicleTheme *model.Theme
	var themeState model.State
	var ok bool

	// @todo find recent state
	fmt.Print("sending default state ", currentState, " ... ")

	err := bgPlayer.SendState(currentState)
	if err != nil {
		fmt.Println("failed")
	} else {
		fmt.Println("ok")
	}

	for {
		loadData(host)
		parseInput(conf)
		// if !input.GameRunning {
		// 	fmt.Println("game not running yet")
		// 	time.Sleep(sleepOffline)
		// 	continue
		// }
		// vehicle changed
		if currentVehicle != input.PlayerVehicle {
			fmt.Print("vehicle change to")
			currentVehicle = input.PlayerVehicle
			if currentVehicle != "" {
				fmt.Println(" " + currentVehicle)
				vehicleConf = getCurrentVehicle()
				// fmt.Println("vehicle", vehicleConf)
				vehicleTheme = conf.GetThemeForVehicle(vehicleConf)
				// fmt.Println("vehicle theme", utils.JsonPretty(vehicleTheme))
				bgPlayer.SendEventStates(&model.EventStates{
					Events: vehicleTheme.Events,
					States: vehicleTheme.States,
				})
				// fmt.Println("sent")
			} else {
				fmt.Println(" none")
			}
		}
		// fast forward state
		newState = ""
		for {
			current_ts = time.Now().Unix()
			// fmt.Println("special loop", newState, current_ts)
			state, err = stMachine.GetNextState(inputMapBool)
			if err != nil {
				fmt.Println("getNextState failed", err)
				cooldown_s = 0
				break
			}
			// fmt.Println("st", state, stMachine.GetCurrentState())
			if state != "" {
				newState = state

				// check for cooldown -1
				if vehicleTheme != nil && cooldown_s > 0 {
					fmt.Println("checking cooldown reset", "state", newState)
					themeState, ok = vehicleTheme.States[newState]
					if ok && themeState.BreaksCooldown == 1 {
						cooldown_s = 0
						fmt.Println("reset cooldown", cooldown_s, "s", "state", newState)
					}
				}
				// state change cooldown, do not set state in colldown period
				if current_ts < (state_ts + cooldown_s) {
					fmt.Println("cooldown in effect", current_ts, state_ts, cooldown_s)
					time.Sleep(time.Millisecond * 100)
					continue
				}

				stMachine.SetState(state)
				state_ts = time.Now().Unix()
				fmt.Println("state", state)
			} else {
				if newState != "" {
					fmt.Println("new state:", newState)
					if vehicleTheme != nil {
						themeState, ok := vehicleTheme.States[newState]
						if ok {
							cooldown_s = themeState.Cooldown
						} else {
							cooldown_s = 0
						}
						fmt.Println("cooldown", cooldown_s, "s")
					}
					bgPlayer.SendState(newState)
					currentState = newState
				} else {
					// fmt.Println("state not changed")
				}
				break
			}
			time.Sleep(time.Millisecond * 100)
		}

		// fmt.Println(utils.JsonPretty(inputMapBool))
		time.Sleep(sleepTime)
	}
}

func Close() {
	// close(chRead)
	// close(chInput)
}
