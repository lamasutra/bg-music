package input

import (
	"math"
	"math/rand/v2"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/wt-client/internal/model"
	"github.com/lamasutra/bg-music/wt-client/internal/player"
)

type inputLoop struct {
	host                          string
	sleepTime                     time.Duration
	conf                          *model.Config
	stMachine                     *model.StateMachine
	bgPlayer                      player.BgPlayer
	currentState                  string
	state                         string
	newState                      string
	currentVehicle                string
	vehicleConf                   *model.Vehicle
	vehicleTheme                  *model.Theme
	hudMsgChecked                 bool
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
	parser                        *DataParser
}

func CreateInputLoop(conf *model.Config, stMachine *model.StateMachine, bgPlayer player.BgPlayer) *inputLoop {
	hmp := model.NewDamageParser()
	return &inputLoop{
		sleepTime:    time.Millisecond * 500,
		host:         conf.Host,
		conf:         conf,
		stMachine:    stMachine,
		bgPlayer:     bgPlayer,
		parser:       createParser(),
		hudMsgParser: hmp,
		player:       hmp.FindOrCreatePlayer(conf.Nickname),
	}
}

func (l *inputLoop) Run() {
	input.Clear()
	// ui.Input(input)

	// @todo find recent state
	logger.Debug("sending default state ", l.currentState, " ... ")

	err := l.bgPlayer.SendState(l.currentState)
	if err != nil {
		logger.Error("failed")
	} else {
		logger.Debug("ok")
	}

	for {
		ok, err := l.checkGameIsUpAndRunning()
		if !ok || err != nil {
			// logger.Debug(err)
			if input.GameRunning {
				logger.Debug("game shut down")
			}
			input.GameRunning = false
			(*inputMapBool)["GameRunning"] = input.GameRunning
			l.handleNextState()
			time.Sleep(l.sleepTime)
			continue
		} else {
			if !input.GameRunning {
				logger.Debug("game is up and running")
			}
			input.GameRunning = true
			(*inputMapBool)["GameRunning"] = input.GameRunning
		}

		l.parser.loadData(l.host)
		if !l.hudMsgChecked {
			l.checkHudMsg()
		}
		l.parser.parseInput(l.conf, l.hudMsgParser, l.player)
		// ui.Input(input)
		// events
		// jstr, _ := json.MarshalIndent(player, "", "  ")
		// logger.Debug(string(jstr))
		if input.MissionStarted {
			l.handleMissionEvents()
		} else {
			l.awacsPermition = false
			l.headingReported = false
			l.headingSpoken = false
			// l.player.Reset()
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
			if !input.PlayerDead && !input.PlayerSeverelyDamaged && !l.player.IsDrone {
				l.bgPlayer.TriggerEvent("damaged")
			}
		}()
	}
	if l.lastPlayerSeverelyDamagedTime < l.player.LastSeverelyDamagedTime {
		l.lastPlayerSeverelyDamagedTime = l.player.LastSeverelyDamagedTime
		go func() {
			r := rand.Float64() * 0.5
			time.Sleep(time.Duration((0.5 + r) * float64(time.Second)))
			if !input.PlayerDead && !l.player.IsDrone {
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
	if !l.player.LastComputedSpeed.IsValid() {
		if !input.PlayerDead {
			l.bgPlayer.Speak("oh no")
		} else {
			// we will try to reset here, probably not the same player
			// l.player.Reset()
		}
	}
}

func (l *inputLoop) handleVehicleChange() {
	l.currentVehicle = input.PlayerVehicle
	if l.currentVehicle != "" {
		logger.Debug("vehicle change to", currentVehicle)
		l.vehicleConf = l.parser.getCurrentVehicle()
		logger.Debug("vehicle", l.vehicleConf)
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
		logger.Debug("vehicle change to none")
		// @todo - send default theme ?
		l.lastKillTime = 0
	}
	logger.Debug("player type:", input.PlayerType)
	if input.EnemyAirNear {
		logger.Debug("air danger")
	} else if input.EnemyAirClose {
		logger.Debug("air combat")
	}
	if input.EnemyGroundNear {
		logger.Debug("ground danger")
	} else if input.EnemyGroundClose {
		logger.Debug("ground combat")
	}
}

func (l *inputLoop) handleNextState() {
	l.newState = ""
	var err error
	for {
		l.state, err = l.stMachine.GetNextState(inputMapBool)
		if err != nil {
			logger.Error("getNextState failed", err)
			break
		}
		// fast forward state
		if l.state != "" {
			l.newState = l.state
			l.stMachine.SetState(l.state)
			logger.Debug("state", l.state)
		} else {
			if l.newState != "" {
				logger.Debug("new state:", l.newState)
				l.bgPlayer.SendState(l.newState)
				l.currentState = l.newState
			} else {
				// fmt.Println("state not changed")
			}
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(l.sleepTime)
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

func (l *inputLoop) checkHudMsg() {
	logger.Debug("checkHudMsg")
	l.hudMsgChecked = true
	l.parser.inputData.HudMsg.Load(l.host, state.lastEvt, state.lastDmg)
	lastDmg := l.parser.inputData.HudMsg.GetLastDmg()
	if lastDmg == nil {
		logger.Debug("not necessary")
		return
	}
	state.lastDmg = uint64(lastDmg.ID)
	logger.Debug("set last id", lastDmg.ID)
	l.parser.inputData.HudMsg.Load(l.host, state.lastEvt, state.lastDmg)
}
