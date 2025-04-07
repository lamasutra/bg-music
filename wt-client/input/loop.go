package input

import (
	"math"
	"math/rand/v2"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/player"
	"github.com/lamasutra/bg-music/wt-client/stateMachine"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

type inputLoop struct {
	pid                           int
	host                          string
	conf                          *model.Config
	stMachine                     *stateMachine.StateMachine
	bgPlayer                      player.BgPlayer
	currentState                  string
	state                         string
	newState                      string
	currentVehicle                string
	vehicleConf                   *model.Vehicle
	vehicleTheme                  *model.Theme
	lastKillTime                  int64
	lastPlayerBurningTime         int64
	lastPlayerDamagedTime         int64
	lastPlayerSeverelyDamagedTime int64
	lastPlayerMadeSeverDamage     int64
	awacsPermition                bool
	headingReported               bool
	headingSpoken                 bool
	hudMsgParser                  *model.DamageParser
	player                        *model.Player
	currentTarget                 *model.Player
}

func CreateInputLoop(conf *model.Config, stMachine *stateMachine.StateMachine, bgPlayer player.BgPlayer) *inputLoop {
	hmp := model.NewDamageParser()
	return &inputLoop{
		host:         conf.Host,
		conf:         conf,
		stMachine:    stMachine,
		bgPlayer:     bgPlayer,
		hudMsgParser: hmp,
		player:       hmp.FindOrCreatePlayer(conf.Nickname),
	}
}

func (l *inputLoop) Run() {
	input.Clear()
	ui.Input(input)

	// @todo find recent state
	ui.Debug("sending default state ", l.currentState, " ... ")

	err := l.bgPlayer.SendState(l.currentState)
	if err != nil {
		ui.Error("failed")
	} else {
		ui.Debug("ok")
	}

	for {
		ok, err := l.checkGameIsUpAndRunning()
		if !ok || err != nil {
			// ui.Debug(err)
			if input.GameRunning {
				ui.Debug("game shut down")
			}
			input.GameRunning = false
			(*inputMapBool)["GameRunning"] = input.GameRunning
			l.handleNextState()
			time.Sleep(sleepTime)
			continue
		} else {
			if !input.GameRunning {
				ui.Debug("game is up and running")
			}
			input.GameRunning = true
			(*inputMapBool)["GameRunning"] = input.GameRunning
		}

		loadData(l.host)
		parseInput(l.conf, l.hudMsgParser, l.player)
		// ui.Input(input)
		// events
		// jstr, _ := json.MarshalIndent(player, "", "  ")
		// ui.Debug(string(jstr))
		if input.MissionStarted {
			l.handleMissionEvents()
		} else {
			l.awacsPermition = false
			l.headingReported = false
			l.headingSpoken = false
			l.player.Reset()
		}
		// vehicle changed
		if l.currentVehicle != input.PlayerVehicle {
			l.handleVehicleChange()
		}

		l.handleNextState()
	}
}

func (l *inputLoop) handleMissionEvents() {
	l.currentTarget = l.player.CurrentTarget
	// @todo - configurable events
	if l.lastKillTime < l.player.LastKillTime {
		l.lastKillTime = l.player.LastKillTime
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((1.5 + r) * float64(time.Second)))
			if !input.PlayerDead {
				l.bgPlayer.TriggerEvent("airKill")
			}
		}()
	}
	if l.lastPlayerBurningTime < l.player.LastBurnedTime {
		l.lastPlayerBurningTime = l.player.LastBurnedTime
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
			if !input.PlayerDead {
				l.bgPlayer.TriggerEvent("burning")
			}
		}()
	}
	if l.lastPlayerDamagedTime < l.player.LastDamagedTime {
		l.lastPlayerDamagedTime = l.player.LastDamagedTime
		input.PlayerDamaged = true
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
			if !input.PlayerDead && !input.PlayerSeverelyDamaged {
				l.bgPlayer.TriggerEvent("damaged")
			}
		}()
	}
	if l.lastPlayerSeverelyDamagedTime < l.player.LastSeverelyDamagedTime {
		l.lastPlayerSeverelyDamagedTime = l.player.LastSeverelyDamagedTime
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
			if !input.PlayerDead {
				l.bgPlayer.TriggerEvent("severely_damaged")
			}
		}()
	}
	if l.lastPlayerMadeSeverDamage < l.player.LastSeverDamageTime {
		l.lastPlayerMadeSeverDamage = l.player.LastSeverDamageTime
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
			if !input.PlayerDead && !l.currentTarget.Dead {
				l.bgPlayer.TriggerEvent("foe_sever_damage")
			}
		}()
	}
	if l.headingSpoken && !l.awacsPermition && input.EnemyAirCount > 1 {
		l.awacsPermition = true
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((10 + r) * float64(time.Second)))
			if !input.PlayerDead {
				l.bgPlayer.TriggerEvent("permitEngage")
			}
		}()
	}
	if !l.headingReported && input.EnemyHeading < 1000 {
		l.headingReported = true
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((3 + r) * float64(time.Second)))
			heading := model.Heading(math.Round(input.EnemyHeading))
			if !input.PlayerDead {
				l.bgPlayer.Speak("hostiles," + strings.Join(heading.Narrate(), ",") + ",degrees")
			}
			l.headingSpoken = true
		}()
	}
}

func (l *inputLoop) handleVehicleChange() {
	l.currentVehicle = input.PlayerVehicle
	if l.currentVehicle != "" {
		ui.Debug("vehicle change to", currentVehicle)
		l.vehicleConf = getCurrentVehicle()
		ui.Debug("vehicle", l.vehicleConf)
		l.vehicleTheme = l.conf.GetThemeForVehicle(l.vehicleConf)
		// fmt.Println("vehicle theme", utils.JsonPretty(vehicleTheme))
		l.bgPlayer.SendEventStates(&model.BgPlayerConfig{
			Events:  l.vehicleTheme.Events,
			States:  l.vehicleTheme.States,
			Narrate: l.vehicleTheme.Narrate,
		})
		l.bgPlayer.ChangeMusic()
		l.awacsPermition = false
		l.headingReported = false
		l.headingSpoken = false
		l.player.Damaged = false
		l.player.SeverlyDamaged = false
		// fmt.Println("sent")
	} else {
		ui.Debug("vehicle change to none")
		// @todo - send default theme ?
		l.lastKillTime = 0
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

func (l *inputLoop) handleNextState() {
	l.newState = ""
	var err error
	for {
		l.state, err = l.stMachine.GetNextState(inputMapBool)
		if err != nil {
			ui.Error("getNextState failed", err)
			break
		}
		// fast forward state
		if l.state != "" {
			l.newState = l.state
			l.stMachine.SetState(l.state)
			ui.Debug("state", state)
		} else {
			if l.newState != "" {
				ui.Debug("new state:", l.newState)
				l.bgPlayer.SendState(l.newState)
				l.currentState = l.newState
			} else {
				// fmt.Println("state not changed")
			}
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(sleepTime)
}

func (l *inputLoop) checkGameIsUpAndRunning() (bool, error) {
	pid, err := l.getPid()

	if err != nil {
		return false, err
	}

	err = syscall.Kill(pid, 0)

	return err == nil, nil
}

func (l *inputLoop) getPid() (int, error) {
	cmd := exec.Command("pidof", "aces")
	out, err := cmd.Output()
	if err != nil {
		return -1, err
	}
	pid, err := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64)
	if err != nil {
		return -1, err
	}

	return int(pid), nil
}
