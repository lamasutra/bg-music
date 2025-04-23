package model

import (
	"time"

	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/wt-client/internal/client"
)

// we have this here to prevent circular dependencies
// I have to admit I am bit lame at this ;)

type WtInput struct {
	GameMode                      string
	PlayerType                    string
	PlayerVehicle                 string
	GameRunning                   bool
	MapLoaded                     bool
	MissionStarted                bool
	MissionEnded                  bool
	PlayerDamaged                 bool
	PlayerSeverelyDamaged         bool
	PlayerDead                    bool
	PlayerLanded                  bool
	EnemyAirNear                  bool
	EnemyGroundNear               bool
	EnemyBaseNear                 bool
	EnemyAirClose                 bool
	EnemyGroundClose              bool
	EnemyBaseClose                bool
	EnemyAirBehind                bool
	IsTanksGameMode               bool
	LastAnyKillTime               int64
	LastPlayerMadeKillTime        int64
	LastPlayerMadeSeverDamageTime int64
	LastPlayerCritDamageTime      int64
	LastPlayerBurningTime         int64
	LastPlayerSeverDamageTime     int64
	NearestEnemyAir               float64
	NearestEnemyGround            float64
	EnemyAirCount                 int
	// temporary
	EnemyHeading float64
}

type WtInputMapBool map[string]bool

type WtData struct {
	Identity     uint32
	HudMsg       *client.HudMsg
	Indicators   *client.Indicators
	MapInfo      *client.MapInfo
	MapObj       *client.MapObj
	State        *client.State
	sleepTime    time.Duration
	sleepOffline time.Duration
	state        struct {
		lastEvt uint64
		lastDmg uint64
	}
}

func NewInputData() *WtData {
	return &WtData{
		sleepTime:    time.Millisecond * 500,
		sleepOffline: time.Millisecond * 1000,
		State:        &client.State{},
		MapInfo:      &client.MapInfo{},
		MapObj:       &client.MapObj{},
		Indicators:   &client.Indicators{},
		HudMsg:       &client.HudMsg{},
	}
}

func (w *WtData) Load(host string) {
	// fmt.Println("Loading data from", host)
	err := w.State.Load(host)
	if err != nil {
		// fmt.Println("state error: ", err)
		time.Sleep(w.sleepOffline)
		err = w.State.Load(host)
		if err != nil {
			return
		}
	}

	err = w.Indicators.Load(host)
	if err != nil {
		logger.Error("indicators error: ", err)
	}
	err = w.MapInfo.Load(host)
	if err != nil {
		logger.Error("!!! mapInfo error: ", err)
	}
	// load other data
	if w.MapInfo.Valid {
		// load map identity
		if w.Identity == 0 {
			w.Identity, err = client.MapIdentity(host)
			logger.Error("map identity error: ", err)
		}

		err = w.MapObj.Load(host)
		if err != nil {
			logger.Error("mapObj error: ", err)
			// } else {
			// fmt.Println(mapObj)
		}
		err = w.HudMsg.Load(host, w.state.lastEvt, w.state.lastDmg)
		if err != nil {
			logger.Error("hudMsg error: ", err)
		} else {

		}
	} else {
		if w.Identity != 0 {
			w.Identity = 0
		}
	}
}

func (w *WtData) SetLastDmg(lastDmg uint64) {
	w.state.lastDmg = lastDmg
}

func (w *WtData) SetLastEvt(lastEvt uint64) {
	w.state.lastEvt = lastEvt
}

func (w *WtData) LoadHudMsg(host string) error {
	return w.HudMsg.Load(host, w.state.lastDmg, w.state.lastEvt)
}

// GameRunning        =     false
// MapLoaded           = false         bool

func (w *WtInput) OnMapNotLoaded() {
	w.MapLoaded = false
	w.GameMode = "unknown"
	w.GameMode = ""
	w.MissionStarted = false
	w.PlayerDamaged = false
	w.PlayerSeverelyDamaged = false
	w.PlayerDead = false
	w.PlayerLanded = false
	w.EnemyAirClose = false
	w.EnemyAirNear = false
	w.EnemyAirCount = 0
	w.EnemyGroundClose = false
	w.EnemyGroundNear = false
	w.IsTanksGameMode = false
	w.PlayerType = ""
	w.PlayerVehicle = ""
	w.LastPlayerMadeKillTime = 0
	w.LastAnyKillTime = 0
	w.LastPlayerBurningTime = 0
	w.LastPlayerCritDamageTime = 0
	w.LastPlayerSeverDamageTime = 0
	w.LastPlayerMadeSeverDamageTime = 0
	w.NearestEnemyAir = -1.0
	w.NearestEnemyGround = -1.0
}

func (w *WtInput) Clear() {
	w.GameMode = "unknown"
	w.PlayerType = ""
	w.PlayerVehicle = ""
	w.MissionStarted = false
	w.MissionEnded = false
	w.PlayerDamaged = false
	w.PlayerDead = false
	w.PlayerLanded = false
	w.EnemyAirNear = false
	w.EnemyGroundNear = false
	w.EnemyBaseNear = false
	w.EnemyAirClose = false
	w.EnemyGroundClose = false
	w.EnemyBaseClose = false
	w.EnemyAirBehind = false
	w.IsTanksGameMode = false
	w.LastAnyKillTime = 0
	w.LastPlayerMadeKillTime = 0
	w.LastPlayerMadeSeverDamageTime = 0
	w.LastPlayerCritDamageTime = 0
	w.LastPlayerBurningTime = 0
	w.LastPlayerSeverDamageTime = 0
	w.NearestEnemyAir = 0
	w.NearestEnemyGround = 0
	w.EnemyAirCount = 0
	// temporary
	w.EnemyHeading = -1000
}

func (w *WtInput) UpdateBoolMap(im *WtInputMapBool, currentTs int64) {
	(*im)["MapLoaded"] = w.MapLoaded
	(*im)["MapLoaded"] = w.MapLoaded
	(*im)["MissionEnded"] = w.MissionEnded
	(*im)["MissionStarted"] = w.MissionStarted
	(*im)["PlayerDamaged"] = w.PlayerDamaged
	(*im)["PlayerSeverelyDamaged"] = w.PlayerSeverelyDamaged
	(*im)["PlayerDead"] = w.PlayerDead
	(*im)["PlayerLanded"] = w.PlayerLanded
	(*im)["EnemyAirClose"] = w.EnemyAirClose
	(*im)["EnemyAirNear"] = w.EnemyAirNear
	(*im)["EnemyGroundClose"] = w.EnemyGroundClose
	(*im)["EnemyGroundNear"] = w.EnemyGroundNear
	(*im)["IsTanksGameMode"] = w.IsTanksGameMode
	(*im)["AirDanger"] = w.EnemyAirNear
	(*im)["AirBattle"] = w.EnemyAirClose || w.LastPlayerMadeKillTime+30 > currentTs || w.LastAnyKillTime+30 > currentTs
	(*im)["PlayerDamaged"] = w.PlayerDamaged
	(*im)["PlayerSeverelyDamaged"] = w.PlayerSeverelyDamaged
	(*im)["GroundDanger"] = w.EnemyGroundNear
	(*im)["GroundBattle"] = w.EnemyGroundClose
}
