package types

import "github.com/lamasutra/bg-music/wt-client/client"

type WtInput struct {
	GameRunning      bool
	MapLoaded        bool
	MissionStarted   bool
	MissionEnded     bool
	PlayerDead       bool
	PlayerLanded     bool
	EnemyAirNear     bool
	EnemyGroundNear  bool
	EnemyBaseNear    bool
	EnemyAirClose    bool
	EnemyGroundClose bool
	EnemyBaseClose   bool
	PlayerType       string
	PlayerVehicle    string
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
