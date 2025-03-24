package input

import (
	"math"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/player"
	"github.com/lamasutra/bg-music/wt-client/stateMachine"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

func LoadLoop(host string, conf *model.Config, stMachine *stateMachine.StateMachine, bgPlayer player.BgPlayer) {
	input.Clear()
	ui.Input(input)
	// fmt.Println("game mode", input.GameMode)
	currentState := stMachine.GetCurrentState()
	var newState, state string
	// var current_ts int64
	var currentVehicle string
	var vehicleConf *model.Vehicle
	var vehicleTheme *model.Theme
	// var themeState model.State
	var lastKillTime int64
	var lastPlayerBurningTime int64
	var lastPlayerDamagedTime int64
	var lastPlayerSeverelyDamagedTime int64
	var lastPlayerMadeSeverDamage int64
	var awacsPermition bool
	var headingReported bool
	var headingSpoken bool
	// var ok bool

	initPatterns(conf)

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
		ui.Input(input)
		// events
		if input.MissionStarted {
			if lastKillTime < input.LastPlayerMadeKillTime {
				lastKillTime = input.LastPlayerMadeKillTime
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((1.5 + r) * float64(time.Second)))
					if !input.PlayerDead {
						bgPlayer.TriggerEvent("airKill")
					}
				}()
			}
			if lastPlayerBurningTime < input.LastPlayerBurningTime {
				lastPlayerBurningTime = input.LastPlayerBurningTime
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
					if !input.PlayerDead {
						bgPlayer.TriggerEvent("burning")
					}
				}()
			}
			if lastPlayerDamagedTime < input.LastPlayerCritDamageTime {
				lastPlayerDamagedTime = input.LastPlayerCritDamageTime
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
					if !input.PlayerDead && !input.PlayerSeverelyDamaged {
						bgPlayer.TriggerEvent("damaged")
					}
				}()
			}
			if lastPlayerSeverelyDamagedTime < input.LastPlayerSeverDamageTime {
				lastPlayerSeverelyDamagedTime = input.LastPlayerSeverDamageTime
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
					if !input.PlayerDead {
						bgPlayer.TriggerEvent("severely_damaged")
					}
				}()
			}
			if lastPlayerMadeSeverDamage < input.LastPlayerMadeSeverDamageTime {
				lastPlayerMadeSeverDamage = input.LastPlayerMadeSeverDamageTime
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
					if !input.PlayerDead {
						bgPlayer.TriggerEvent("foe_sever_damage")
					}
				}()
			}
			if headingSpoken && !awacsPermition && input.EnemyAirCount > 1 {
				awacsPermition = true
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((10 + r) * float64(time.Second)))
					if !input.PlayerDead {
						bgPlayer.TriggerEvent("permitEngage")
					}
				}()
			}
			if !headingReported && input.EnemyHeading < 1000 {
				headingReported = true
				go func() {
					r := rand.Float64() * 0.5
					time.Sleep(time.Duration((3 + r) * float64(time.Second)))
					heading := model.Heading(math.Round(input.EnemyHeading))
					if !input.PlayerDead {
						bgPlayer.Speak("hostiles," + strings.Join(heading.Narrate(), ",") + ",degrees")
					}
					headingSpoken = true
				}()
			}
		} else {
			awacsPermition = false
			headingReported = false
			headingSpoken = false
		}
		// vehicle changed
		if currentVehicle != input.PlayerVehicle {
			currentVehicle = input.PlayerVehicle
			if currentVehicle != "" {
				ui.Debug("vehicle change to", currentVehicle)
				vehicleConf = getCurrentVehicle()
				ui.Debug("vehicle", vehicleConf)
				vehicleTheme = conf.GetThemeForVehicle(vehicleConf)
				// fmt.Println("vehicle theme", utils.JsonPretty(vehicleTheme))
				bgPlayer.SendEventStates(&model.BgPlayerConfig{
					Events:  vehicleTheme.Events,
					States:  vehicleTheme.States,
					Narrate: vehicleTheme.Narrate,
				})
				awacsPermition = false
				headingReported = false
				headingSpoken = false
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
			state, err = stMachine.GetNextState(inputMapBool)
			if err != nil {
				ui.Error("getNextState failed", err)
				break
			}
			if state != "" {
				newState = state
				stMachine.SetState(state)
				ui.Debug("state", state)
			} else {
				if newState != "" {
					ui.Debug("new state:", newState)
					bgPlayer.SendState(newState)
					currentState = newState
				} else {
					// fmt.Println("state not changed")
				}
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		time.Sleep(sleepTime)
	}
}
