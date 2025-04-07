package types

import "github.com/lamasutra/bg-music/wt-client/client"

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
	Identity   uint32
	HudMsg     *client.HudMsg
	Indicators *client.Indicators
	MapInfo    *client.MapInfo
	MapObj     *client.MapObj
	State      *client.State
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
