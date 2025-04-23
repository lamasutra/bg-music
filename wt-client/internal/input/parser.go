package input

import (
	"math"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/wt-client/internal/client"
	"github.com/lamasutra/bg-music/wt-client/internal/model"
)

func createParser() *DataParser {
	gameModes := make(map[string]GameMode, 3)
	gameModes["default"] = NewDefaultMode()
	gameModes["tanks"] = NewTanksMode()
	gameModes["air"] = NewAirMode()

	return &DataParser{
		modes:       gameModes,
		currentMode: gameModes["default"],
		inputData:   model.NewInputData(),
	}
}

func (p *DataParser) loadData(host string) {
	p.inputData.Load(host)
}

func (p *DataParser) parseInput(conf *model.Config, hudMsgParser *model.DamageParser, player *model.Player) {
	// fmt.Println("Data:", inputData)

	input.MapLoaded = p.inputData.MapInfo.Valid
	var objDistance float64
	var nearestAirfield *client.Entity
	var nearestTankRespawnBase *client.Entity
	var tankRespawnBases []client.Entity

	current_ts := time.Now().Unix()

	if input.MapLoaded {
		hudMsgParser.Parse(p.inputData.HudMsg)

		// fmt.Println("game mode", input.GameMode)
		playerEntity := p.inputData.MapObj.GetPlayerEntity()
		airfields := p.inputData.MapObj.GetAirfields()
		tankRespawnBases = *p.inputData.MapObj.GetTankRespawnBases()
		// if len(*airfields) > 0 {
		// fmt.Println("we have airfields")
		if playerEntity != nil {
			nearestAirfield = p.getNearestEntity(playerEntity, airfields, p.inputData.MapObj, p.inputData.MapInfo)
			if input.GameMode == "tanks" && len(tankRespawnBases) > 0 {
				nearestTankRespawnBase = p.getNearestEntity(playerEntity, &tankRespawnBases, p.inputData.MapObj, p.inputData.MapInfo)
			}
		}
		// } else {
		// fmt.Println("airfields not found")
		// }
		if nearestAirfield != nil && p.inputData.State != nil {
			player.LoadData(playerEntity, nearestAirfield, p.inputData.State, p.inputData.Indicators, p.inputData.MapInfo, p.inputData.MapObj)
		}
		if tankRespawnBases != nil && nearestTankRespawnBase != nil {
			player.CheckIsSpawned(nearestTankRespawnBase, p.inputData.MapInfo, p.inputData.MapObj)
		}
		enemyAircrafts := p.inputData.MapObj.GetAircraftsByColors(&conf.Colors.Foe.Air)
		enemyGroundUnits := p.inputData.MapObj.GetGroundUnitsByColors(&conf.Colors.Foe.Ground)
		captureZones := p.inputData.MapObj.GetCaptureZones()
		if input.GameMode == "unknown" {
			if len(tankRespawnBases) > 0 {
				input.GameMode = "tanks"
				input.IsTanksGameMode = true
			} else {
				input.GameMode = "air"
				input.IsTanksGameMode = false
			}
			logger.Debug("Game mode detected:", input.GameMode)
			// reset last dmg id
			lastDamage := p.inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				logger.Debug("reseting lastDmg to", lastDamage.ID)
				p.inputData.SetLastDmg(uint64(lastDamage.ID))
			}
		}
		input.EnemyAirCount = len(*enemyAircrafts)
		input.MissionEnded = false
		if !input.MissionStarted {
			// set lastDamage to the latest on mission start
			lastDamage := p.inputData.HudMsg.GetLastDmg()
			if lastDamage != nil {
				p.inputData.SetLastDmg(uint64(lastDamage.ID))
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
				input.PlayerDamaged = true
				input.LastPlayerCritDamageTime = lastCritTime
			}
			lastSeverTime := player.LastSeverelyDamagedTime
			if lastSeverTime > 0 {
				input.LastPlayerSeverDamageTime = lastSeverTime
				input.PlayerDamaged = true
				input.PlayerSeverelyDamaged = true
			}
		}
		input.PlayerDead = input.MissionStarted && (player.Dead || shouldStayDead)
		// input.PlayerLanded = input.MissionStarted && playerEntity == nil && !player.Dead
		input.PlayerLanded = player.Landed && !player.Dead
		if player.Dead && input.GameMode == "air" {
			logger.Debug("air, player should stay dead")
			shouldStayDead = true
		}
		// @todo add speed, landed is alive and 0 speed
		if input.PlayerLanded {
			logger.Debug("player landed, reseting damage")
			player.Damaged = false
			player.SeverlyDamaged = false
		}
		// reset if spawned
		if player.Spawned {
			player.Damaged = false
			player.SeverlyDamaged = false
			player.Dead = false
		}
		// fmt.Println("dead conds", input.PlayerDead, isShotDown, hasCrashed)

		// vehicle changed
		if input.PlayerVehicle != p.inputData.Indicators.Type {
			input.PlayerType = p.inputData.Indicators.Army
			input.PlayerVehicle = p.inputData.Indicators.Type
			currentVehicle = conf.GetVehicleForPlayerTypeAndVehicleType(input.PlayerType, input.PlayerVehicle)
			currentTheme = conf.GetThemeForVehicle(currentVehicle)

			// reset if vehicle changed
			player.Damaged = false
			player.SeverlyDamaged = false
			player.Dead = false
		}

		// @todo last known location
		if playerEntity != nil {
			// reset player damage and dead state if in tank mode, reset after 10s, we are unable to process state of the player from map objects
			// if player.Dead && input.GameMode == "tanks" && player.LastKilledTime+10 > current_ts {
			// 	player.Dead = false
			// 	player.Damaged = false
			// }
			nearestAir := p.getNearestEntity(playerEntity, enemyAircrafts, p.inputData.MapObj, p.inputData.MapInfo)
			nearestGround := p.getNearestEntity(playerEntity, enemyGroundUnits, p.inputData.MapObj, p.inputData.MapInfo)

			if nearestAir != nil {
				objDistance = p.inputData.MapObj.GetDistance(playerEntity, nearestAir, p.inputData.MapInfo)
				input.NearestEnemyAir = objDistance
				// 10000
				input.EnemyAirNear = objDistance < float64(currentTheme.Distances.Air.Danger)
				// 5000
				input.EnemyAirClose = objDistance < float64(currentTheme.Distances.Air.Combat)
				input.EnemyHeading = p.inputData.MapObj.GetHeading(playerEntity, nearestAir)
				// fmt.Println("air dist", objDistance, float64(currentTheme.Distances.Air.Danger), float64(currentTheme.Distances.Air.Combat), input.EnemyAirNear, input.EnemyAirClose)
			} else {
				input.EnemyAirNear = false
				input.EnemyAirClose = false
				input.NearestEnemyAir = -1.0
				input.EnemyHeading = 1000
			}

			if nearestGround != nil {
				objDistance = p.inputData.MapObj.GetDistance(playerEntity, nearestGround, p.inputData.MapInfo)
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

			if input.GameMode == "tanks" {
				nearestZone := p.getNearestEntity(playerEntity, captureZones, p.inputData.MapObj, p.inputData.MapInfo)
				objDistance = p.inputData.MapObj.GetDistance(playerEntity, nearestZone, p.inputData.MapInfo)
				if !input.EnemyGroundNear {
					input.EnemyGroundNear = objDistance < float64(currentTheme.Distances.Ground.Danger)
				}
				if !input.EnemyGroundClose {
					input.EnemyGroundClose = objDistance < float64(currentTheme.Distances.Ground.Combat)
				}
			}

			// @todo enemy base
			// } else {
			// input.EnemyAirNear = false
			// input.EnemyAirClose = false
			// if input.PlayerLanded {
			// input.PlayerDamaged = false
			// input.PlayerSeverelyDamaged = false
			// }
		}

		// only set lastDamage if player is dead and resurrected
		if p.isMissionEnded(p.inputData.HudMsg) {
			input.MissionEnded = true
		}
		// if input.PlayerDead && player != nil {
		lastDamage := p.inputData.HudMsg.GetLastDmg()
		if lastDamage != nil {
			p.inputData.SetLastDmg(uint64(lastDamage.ID))
		}
		// }
	} else {
		if input.GameMode != "unknown" {
			logger.Debug("no map loaded, reseting current vechicle, should stay dead, game mode")
			input.OnMapNotLoaded()
			currentVehicle = nil
			shouldStayDead = false
			input.GameMode = "unknown"
			input.MissionEnded = true
			player.Reset()
		}
	}

	input.UpdateBoolMap(inputMapBool, current_ts)
}

func (p *DataParser) getCurrentVehicle() *model.Vehicle {
	return currentVehicle
}

func (p *DataParser) getNearestEntity(player *client.Entity, entities *[]client.Entity, mapObj *client.MapObj, mapInfo *client.MapInfo) *client.Entity {
	var nearest *client.Entity
	var distance, lastDistance float64
	lastDistance = math.MaxInt64

	if len(*entities) == 0 {
		return nil
	}
	for _, entity := range *entities {
		distance = mapObj.GetDistance(player, &entity, mapInfo)
		if distance < lastDistance {
			lastDistance = distance
			nearest = &entity
		}
	}

	return nearest
}

func (p *DataParser) isMissionEnded(hudMsg *client.HudMsg) bool {
	for _, msg := range hudMsg.Damage {
		if strings.Contains(msg.Msg, "has delivered the final blow!") {
			return true
		}
	}

	return false
}
