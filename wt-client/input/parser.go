package input

import (
	"regexp"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/model"
)

func parseInput(conf *model.Config) {
	// fmt.Println("Data:", inputData)

	input.MapLoaded = inputData.MapInfo.Valid
	var objDistance float64

	current_ts := time.Now().Unix()

	if input.MapLoaded {
		// fmt.Println("game mode", input.GameMode)
		player := inputData.MapObj.GetPlayerEntity()
		enemyAircrafts := inputData.MapObj.GetAircraftsByColors(&conf.Colors.Foe.Air)
		enemyGroundUnits := inputData.MapObj.GetGroundUnitsByColors(&conf.Colors.Foe.Ground)
		if input.GameMode == "unknown" {
			tankRespawnBases := inputData.MapObj.GetTankRespawnBases()
			if len(*tankRespawnBases) > 0 {
				input.GameMode = "tanks"
			} else {
				input.GameMode = "air"
			}
		}
		input.EnemyAirCount = len(*enemyAircrafts)
		isShotDown := playerIsShotDown(inputData.HudMsg)
		hasCrashed := playerHasCrashed(inputData.HudMsg)
		// isMissionEnded := inputData.hudMsg.IsMissionEnded()
		input.MissionEnded = false
		if !input.MissionStarted {
			// set lastDamage to the latest on mission start
			lastDamage := inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				state.lastDmg = uint64(lastDamage.ID)
			}
			input.MissionStarted = player != nil // @todo what else ?
		} else {
			// parse messages from hud_msg to handle events later
			lastKillMadeTime := getLastPlayerMadeKillTime(inputData.HudMsg)
			if lastKillMadeTime > 0 {
				input.LastPlayerMadeKillTime = lastKillMadeTime
			}
			lastSeverDamageMadeTime := getLastPlayerMadeSeverDamageTime(inputData.HudMsg)
			if lastSeverDamageMadeTime > 0 {
				input.LastPlayerMadeSeverDamageTime = lastSeverDamageMadeTime
			}
			lastAnyKillTime := getLastAnyKillTime(inputData.HudMsg)
			if lastAnyKillTime > 0 {
				input.LastAnyKillTime = lastKillMadeTime
			}
			lastBurningTime := getLastPlayerIsBurningTime(inputData.HudMsg)
			if lastBurningTime > 0 {
				input.LastPlayerBurningTime = lastBurningTime
			}
			lastCritTime := getLastPlayerIsCritDamagedTime(inputData.HudMsg)
			if lastCritTime > 0 {
				input.LastPlayerCritDamageTime = lastCritTime
			}
			lastSeverTime := getLastPlayerIsSeverelyDamagedTime(inputData.HudMsg)
			if lastSeverTime > 0 {
				input.LastPlayerSeverDamageTime = lastSeverTime
			}
		}
		input.PlayerDead = input.MissionStarted && ((isShotDown || hasCrashed) || shouldStayDead)
		input.PlayerLanded = input.MissionStarted && player == nil && !isShotDown && !hasCrashed
		if input.PlayerDead && input.GameMode == "air" {
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
		if player != nil {
			nearestAir := getNearestEntity(player, enemyAircrafts, inputData.MapObj, inputData.MapInfo)
			nearestGround := getNearestEntity(player, enemyGroundUnits, inputData.MapObj, inputData.MapInfo)

			if nearestAir != nil {
				objDistance = inputData.MapObj.GetDistance(player, nearestAir, inputData.MapInfo)
				input.NearestEnemyAir = objDistance
				// 10000
				input.EnemyAirNear = objDistance < float64(currentTheme.Distances.Air.Danger)
				// 5000
				input.EnemyAirClose = objDistance < float64(currentTheme.Distances.Air.Combat)
				input.EnemyHeading = inputData.MapObj.GetHeading(player, nearestAir)
				// fmt.Println("air dist", objDistance, float64(currentTheme.Distances.Air.Danger), float64(currentTheme.Distances.Air.Combat), input.EnemyAirNear, input.EnemyAirClose)
			} else {
				input.EnemyAirNear = false
				input.EnemyAirClose = false
				input.NearestEnemyAir = -1.0
				input.EnemyHeading = 1000
			}

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
		if input.PlayerType != "" {
			input.MissionEnded = true
		}

		input.OnMapNotLoaded()
		currentVehicle = nil
		shouldStayDead = false
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

func playerHasCrashed(hudMsg *client.HudMsg) bool {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerHasCrashedRegExp) > 0
}

func playerIsShotDown(hudMsg *client.HudMsg) bool {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerIsShotDownRegExp) > 0
}

func getLastAnyKillTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &anyKillRegExp)
}

func getLastPlayerMadeKillTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerMadeKilledRegExp)
}

func getLastPlayerMadeSeverDamageTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerMadeSeverDamagedRegExp)
}

func getLastPlayerIsBurningTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerIsBurningRegExp)
}

func getLastPlayerIsCritDamagedTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerIsCritDamagedRegExp)
}

func getLastPlayerIsSeverelyDamagedTime(hudMsg *client.HudMsg) int64 {
	return getLastTimeForHudMsgAndPattern(hudMsg, &playerIsSeverlyDamagedRegExp)
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
