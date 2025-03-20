package input

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/player"
	"github.com/lamasutra/bg-music/wt-client/stateMachine"
	"github.com/lamasutra/bg-music/wt-client/types"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

const sleepTime = time.Millisecond * 500
const sleepOffline = time.Millisecond * 1000

const anyKillPattern = `shot\s+down|destroyed`

var anyKillRegExp regexp.Regexp

const playerShotDownPatternTemplate = `%s.+(shot\s+down|destroyed).+`

var playerShotDownRegExp regexp.Regexp

const playerIsShotDownPatternTemplate = `shot down %s`

var playerIsShotDownRegExp regexp.Regexp

const playerHasCrashedPatternTemplate = `%s.+has crashed`

var playerHasCrashedRegExp regexp.Regexp

// const playerShotDownRegExpTemplate = `%s\s+\([^\)]+\)\s+(shot\s+down|destroyed)\s+` // [^\(]+\([^\)]+\)

// const playerShootDownEnemyRegexp = `[^\(]+\s+\([^\)]+\)\s+shot\s+down\s+[^\(]+\([^\)]+\)`

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
var input = &types.WtInput{
	NearestEnemyAir:    -1,
	NearestEnemyGround: -1,
}
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
		ui.Error("indicators error: ", err)
	}
	err = inputData.MapInfo.Load(host)
	if err != nil {
		ui.Error("mapInfo error: ", err)
	}
	// load other data
	if inputData.MapInfo.Valid {
		// load map identity
		if inputData.Identity == 0 {
			inputData.Identity, err = client.MapIdentity(host)
			ui.Error("map identity error: ", err)
		}

		err = inputData.MapObj.Load(host)
		if err != nil {
			ui.Error("mapObj error: ", err)
			// } else {
			// fmt.Println(mapObj)
		}
		err = inputData.HudMsg.Load(host, lastEvt, lastDmg)
		if err != nil {
			ui.Error("hudMsg error: ", err)
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

	current_ts := time.Now().Unix()

	if input.MapLoaded {
		player := inputData.MapObj.GetPlayerEntity()
		enemyAircrafts := inputData.MapObj.GetEnemyAircrafts(&conf.Colors.Enemy.Air)
		enemyGroundUnits := inputData.MapObj.GetEnemyGroundUnits(&conf.Colors.Enemy.Ground)
		isShotDown := playerIsShotDown(inputData.HudMsg)
		hasCrashed := playerHasCrashed(inputData.HudMsg)
		// isMissionEnded := inputData.hudMsg.IsMissionEnded()
		input.MissionEnded = false
		if !input.MissionStarted {
			// set lastDamage to the latest on mission start
			lastDamage := inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				lastDmg = uint64(lastDamage.ID)
			}
			input.MissionStarted = player != nil // @todo what else ?
		} else {
			lastKillTime := getLastKillTime(inputData.HudMsg)
			if lastKillTime > 0 {
				input.LastKillTime = lastKillTime
			}
			lastAnyKillTime := getLastAnyKillTime(inputData.HudMsg)
			if lastAnyKillTime > 0 {
				input.LastAnyKillTime = lastKillTime
			}
		}
		input.PlayerLanded = input.MissionStarted && player == nil && !isShotDown && !hasCrashed
		input.PlayerDead = input.MissionStarted && player == nil && (isShotDown || hasCrashed)
		fmt.Println("dead conds", input.PlayerDead, isShotDown, hasCrashed)

		// vehicle changed
		if input.PlayerVehicle != inputData.Indicators.Type {
			input.PlayerType = inputData.Indicators.Army
			input.PlayerVehicle = inputData.Indicators.Type
			currentVehicle = conf.GetVehicleForPlayerTypeAndVehicleType(input.PlayerType, input.PlayerVehicle)
			currentTheme = conf.GetThemeForVehicle(currentVehicle)
		}

		// @todo last known location
		if player != nil {
			nearestAir := getNearestEntity(player, enemyAircrafts, inputData.MapObj, inputData.MapInfo)
			nearestGround := getNearestEntity(player, enemyGroundUnits, inputData.MapObj, inputData.MapInfo)

			// @todo configurable in theme
			if nearestAir != nil {
				objDistance = inputData.MapObj.GetDistance(player, nearestAir, inputData.MapInfo)
				input.NearestEnemyAir = objDistance
				// 10000
				input.EnemyAirNear = objDistance < float64(currentTheme.Distances.Air.Danger)
				// 5000
				input.EnemyAirClose = objDistance < float64(currentTheme.Distances.Air.Combat)
				// fmt.Println("air dist", objDistance, float64(currentTheme.Distances.Air.Danger), float64(currentTheme.Distances.Air.Combat), input.EnemyAirNear, input.EnemyAirClose)
			} else {
				input.EnemyAirNear = false
				input.EnemyAirClose = false
				input.NearestEnemyAir = -1.0
			}

			// @todo configurable in theme
			if nearestGround != nil {
				objDistance = inputData.MapObj.GetDistance(player, nearestGround, inputData.MapInfo)
				input.NearestEnemyGround = objDistance
				// 20000
				input.EnemyGroundNear = objDistance < float64(currentTheme.Distances.Ground.Danger)
				// 1000
				input.EnemyGroundClose = objDistance < float64(currentTheme.Distances.Ground.Combat)
			} else {
				// fmt.Println("ground dist", objDistance, float64(currentTheme.Distances.Ground.Danger), float64(currentTheme.Distances.Ground.Combat))
				input.EnemyGroundNear = false
				input.EnemyGroundClose = false
				input.NearestEnemyGround = -1.0
			}

			// @todo enemy base
		} else {
			input.EnemyAirNear = false
			input.EnemyAirClose = false
		}

		// only set lastDamage if player is dead and resurrected
		if isMissionEnded(inputData.HudMsg) {
			input.MissionEnded = true
		}
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
		input.LastKillTime = 0
		input.NearestEnemyAir = -1
		input.NearestEnemyGround = -1
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
	(*inputMapBool)["AirDanger"] = input.EnemyAirNear
	(*inputMapBool)["AirBattle"] = input.EnemyAirClose || input.LastKillTime+30 > current_ts || input.LastAnyKillTime+30 > current_ts
	(*inputMapBool)["GroundDanger"] = input.EnemyGroundNear
	(*inputMapBool)["GroundBattle"] = input.EnemyGroundClose

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
	ui.Input(*input)

	currentState := stMachine.GetCurrentState()
	var newState, state string
	// var current_ts int64
	var cooldown_s int64
	var currentVehicle string
	var vehicleConf *model.Vehicle
	var vehicleTheme *model.Theme
	// var themeState model.State
	var lastKillTime int64
	// var ok bool

	playerShotDownRegExp = *regexp.MustCompile(fmt.Sprintf(playerShotDownPatternTemplate, conf.Nickname))
	anyKillRegExp = *regexp.MustCompile(anyKillPattern)
	playerHasCrashedRegExp = *regexp.MustCompile(fmt.Sprintf(playerHasCrashedPatternTemplate, conf.Nickname))
	playerIsShotDownRegExp = *regexp.MustCompile(fmt.Sprintf(playerIsShotDownPatternTemplate, conf.Nickname))

	fmt.Println(playerShotDownRegExp, anyKillRegExp, playerHasCrashedRegExp, playerIsShotDownRegExp)
	// @todo find recent state
	ui.Debug("sending default state ", currentState, " ... ")

	err := bgPlayer.SendState(currentState)
	if err != nil {
		ui.Error("failed")
	} else {
		ui.Debug("ok")
	}

	for {
		loadData(host)
		parseInput(conf)
		ui.Input(*input)
		if input.MissionStarted && lastKillTime < input.LastKillTime {
			lastKillTime = input.LastKillTime
			bgPlayer.TriggerEvent("airKill")
		}
		// if !input.GameRunning {
		// 	fmt.Println("game not running yet")
		// 	time.Sleep(sleepOffline)
		// 	continue
		// }
		// vehicle changed
		if currentVehicle != input.PlayerVehicle {
			currentVehicle = input.PlayerVehicle
			if currentVehicle != "" {
				ui.Debug("vehicle change to", currentVehicle)
				vehicleConf = getCurrentVehicle()
				fmt.Println("vehicle", vehicleConf)
				vehicleTheme = conf.GetThemeForVehicle(vehicleConf)
				// fmt.Println("vehicle theme", utils.JsonPretty(vehicleTheme))
				bgPlayer.SendEventStates(&model.EventStates{
					Events: vehicleTheme.Events,
					States: vehicleTheme.States,
				})
				// fmt.Println("sent")
			} else {
				ui.Debug("vehicle change to none")
				// @todo - send default theme ?
				lastKillTime = 0
			}
			ui.Debug("player type:", input.PlayerType)
			if input.EnemyAirNear {
				ui.Debug("air danger")
			} else if input.EnemyAirClose {
				ui.Debug("air combat")
			}
			if input.EnemyGroundNear {
				ui.Debug("ground danger")
			} else if input.EnemyGroundClose {
				ui.Debug("ground combat")
			}
		}
		// fast forward state
		newState = ""
		for {
			// current_ts = time.Now().Unix()
			// fmt.Println("special loop", newState, current_ts)
			state, err = stMachine.GetNextState(inputMapBool)
			if err != nil {
				ui.Error("getNextState failed", err)
				// cooldown_s = 0
				break
			}
			// fmt.Println("st", state, stMachine.GetCurrentState())
			if state != "" {
				if state != newState {
					state_ts = time.Now().Unix()
				}
				newState = state

				// check for cooldown -1
				// if vehicleTheme != nil && cooldown_s > 0 {
				// 	ui.Debug("checking cooldown reset", "state", newState)
				// 	themeState, ok = vehicleTheme.States[newState]
				// 	if ok && themeState.BreaksCooldown == 1 {
				// 		cooldown_s = 0
				// 		ui.Debug("reset cooldown", cooldown_s, "s", "state", newState)
				// 	}
				// }
				// state change cooldown, do not set state in colldown period
				// if current_ts < (state_ts + cooldown_s) {
				// 	ui.Debug("cooldown in effect", current_ts, state_ts, cooldown_s)
				// 	time.Sleep(time.Millisecond * 100)
				// 	continue
				// }

				stMachine.SetState(state)
				state_ts = time.Now().Unix()
				ui.Debug("state", state)
			} else {
				if newState != "" {
					ui.Debug("new state:", newState)
					if vehicleTheme != nil {
						themeState, ok := vehicleTheme.States[newState]
						if ok {
							cooldown_s = themeState.Cooldown
						} else {
							cooldown_s = 0
						}
						ui.Debug("cooldown", cooldown_s, "s")
					}
					bgPlayer.SendState(newState)
					currentState = newState
				} else {
					// fmt.Println("state not changed")
				}
				break
			}
			// ui.Debug("inner loop")
			time.Sleep(time.Millisecond * 100)
		}
		// ui.Debug("outer loop")
		time.Sleep(sleepTime)
	}
}

func isMissionEnded(hudMsg *client.HudMsg) bool {
	for _, msg := range hudMsg.Damage {
		if strings.Contains(msg.Msg, "has delivered the final blow!") {
			return true
		}
	}

	return false
}

func playerHasCrashed(hudMsg *client.HudMsg) bool {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerHasCrashedRegExp) > 0
}

func playerIsShotDown(hudMsg *client.HudMsg) bool {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerIsShotDownRegExp) > 0
}

func getLastAnyKillTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &anyKillRegExp)
}

func getLastKillTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerShotDownRegExp)
}

func getLastTimeForHudMsgAndPattern(hudMsg *client.HudMsg, regExp *regexp.Regexp) int64 {
	messages := hudMsg.MatchMessages(regExp)

	// fmt.Println(messages)
	length := len(messages)
	if length == 0 {
		return 0
	}
	lastMsg := messages[length-1]

	return int64(lastMsg.Time)
}

// func loopInput() {

// }

// func loopState() {

// }

func Close() {
	// close(chRead)
	// close(chInput)
}
