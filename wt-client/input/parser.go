package input

import (
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

func parseInput(conf *model.Config, hudMsgParser *model.DamageParser, player *model.Player) {
	// fmt.Println("Data:", inputData)

	input.MapLoaded = inputData.MapInfo.Valid
	var objDistance float64

	current_ts := time.Now().Unix()

	if input.MapLoaded {
		hudMsgParser.Parse(inputData.HudMsg)

		// fmt.Println("game mode", input.GameMode)
		playerEntity := inputData.MapObj.GetPlayerEntity()
		enemyAircrafts := inputData.MapObj.GetAircraftsByColors(&conf.Colors.Foe.Air)
		enemyGroundUnits := inputData.MapObj.GetGroundUnitsByColors(&conf.Colors.Foe.Ground)
		if input.GameMode == "unknown" {
			tankRespawnBases := inputData.MapObj.GetTankRespawnBases()
			if len(*tankRespawnBases) > 0 {
				input.GameMode = "tanks"
			} else {
				input.GameMode = "air"
			}
			ui.Debug("Game mode detected:", input.GameMode)
			// reset last dmg id
			lastDamage := inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				ui.Debug("reseting lastDmg to", lastDamage.ID)
				state.lastDmg = uint64(lastDamage.ID)
			}
		}
		input.EnemyAirCount = len(*enemyAircrafts)
		input.MissionEnded = false
		if !input.MissionStarted {
			// set lastDamage to the latest on mission start
			lastDamage := inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				state.lastDmg = uint64(lastDamage.ID)
			}
			input.MissionStarted = playerEntity != nil // @todo what else ?
		} else {
			// parse messages from hud_msg to handle events later
			lastKillMadeTime := player.LastKillTime
			if lastKillMadeTime > 0 {
				input.LastPlayerMadeKillTime = lastKillMadeTime
			}
			lastSeverDamageMadeTime := player.LastSeverDamageTime
			if lastSeverDamageMadeTime > 0 {
				input.LastPlayerMadeSeverDamageTime = lastSeverDamageMadeTime
			}
			lastAnyKillTime := hudMsgParser.GetLastKillTime()
			if lastAnyKillTime > 0 {
				input.LastAnyKillTime = lastKillMadeTime
			}
			lastBurningTime := player.LastBurnedTime
			if lastBurningTime > 0 {
				input.LastPlayerBurningTime = lastBurningTime
			}
			lastCritTime := player.LastDamagedTime
			if lastCritTime > 0 {
				input.LastPlayerCritDamageTime = lastCritTime
			}
			lastSeverTime := player.LastSeverelyDamagedTime
			if lastSeverTime > 0 {
				input.LastPlayerSeverDamageTime = lastSeverTime
			}
		}
		input.PlayerDead = input.MissionStarted && (player.Dead || shouldStayDead)
		input.PlayerLanded = input.MissionStarted && playerEntity == nil && !player.Dead
		if player.Dead && input.GameMode == "air" {
			shouldStayDead = true
		}
		// fmt.Println("dead conds", input.PlayerDead, isShotDown, hasCrashed)

		// vehicle changed
		if input.PlayerVehicle != inputData.Indicators.Type {
			input.PlayerType = inputData.Indicators.Army
			input.PlayerVehicle = inputData.Indicators.Type
			currentVehicle = conf.GetVehicleForPlayerTypeAndVehicleType(input.PlayerType, input.PlayerVehicle)
			currentTheme = conf.GetThemeForVehicle(currentVehicle)
		}

		// @todo last known location
		if playerEntity != nil {
			// reset player damage and dead state if in tank mode, reset after 10s, we are unable to process state of the player from map objects
			if player.Dead && input.GameMode == "tanks" && player.LastKilledTime+10 > current_ts {
				player.Dead = false
				player.Damaged = false
			}
			nearestAir := getNearestEntity(playerEntity, enemyAircrafts, inputData.MapObj, inputData.MapInfo)
			nearestGround := getNearestEntity(playerEntity, enemyGroundUnits, inputData.MapObj, inputData.MapInfo)

			if nearestAir != nil {
				objDistance = inputData.MapObj.GetDistance(playerEntity, nearestAir, inputData.MapInfo)
				input.NearestEnemyAir = objDistance
				// 10000
				input.EnemyAirNear = objDistance < float64(currentTheme.Distances.Air.Danger)
				// 5000
				input.EnemyAirClose = objDistance < float64(currentTheme.Distances.Air.Combat)
				input.EnemyHeading = inputData.MapObj.GetHeading(playerEntity, nearestAir)
				// fmt.Println("air dist", objDistance, float64(currentTheme.Distances.Air.Danger), float64(currentTheme.Distances.Air.Combat), input.EnemyAirNear, input.EnemyAirClose)
			} else {
				input.EnemyAirNear = false
				input.EnemyAirClose = false
				input.NearestEnemyAir = -1.0
				input.EnemyHeading = 1000
			}

			if nearestGround != nil {
				objDistance = inputData.MapObj.GetDistance(playerEntity, nearestGround, inputData.MapInfo)
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
			if input.PlayerLanded {
				input.PlayerDamaged = false
				input.PlayerSeverelyDamaged = false
			}
		}

		// only set lastDamage if player is dead and resurrected
		if isMissionEnded(inputData.HudMsg) {
			input.MissionEnded = true
		}
		// if input.PlayerDead && player != nil {
		lastDamage := inputData.HudMsg.GetLastDmg()
		if lastDamage != nil {
			state.lastDmg = uint64(lastDamage.ID)
		}
		// }
	} else {
		if input.GameMode != "unknown" {
			ui.Debug("no map loaded, reseting current vechicle, should stay dead, game mode")
			input.OnMapNotLoaded()
			currentVehicle = nil
			shouldStayDead = false
			input.GameMode = "unknown"
			input.MissionEnded = true
		}
	}

	input.UpdateBoolMap(inputMapBool, current_ts)

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

func isMissionEnded(hudMsg *client.HudMsg) bool {
	for _, msg := range hudMsg.Damage {
		if strings.Contains(msg.Msg, "has delivered the final blow!") {
			return true
		}
	}

	return false
}
