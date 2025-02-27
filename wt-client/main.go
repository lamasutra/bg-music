package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/clientConfig"
	"github.com/lamasutra/bg-music/wt-client/player"
)

func writeToFile(value string, file *os.File) {
	_, err := file.WriteString(value + "\n")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var conf clientConfig.Config
	err := conf.Read("wt-config.json")
	if err != nil {
		fmt.Println("Cannot read wt-config.json")
		return
	}

	fmt.Println("your nickname", conf.Nickname)

	player := player.CreatePlayer(conf.BgPlayerType, &conf)

	defer player.Close()
	for {
		time.Sleep(time.Second)
		fmt.Println("tick")
	}
}

func mainBak() {
	var CurrentConfig clientConfig.Config
	var state client.State
	var indicators client.Indicators
	var mapInfo client.MapInfo
	var mapObj client.MapObj
	var hudMsg client.HudMsg
	var lastEvt, lastDmg uint64
	var err error

	err = CurrentConfig.Read("wt-config.json")
	if err != nil {
		fmt.Println("Cannot read wt-config.json")
	}

	fmt.Println("your nickname", CurrentConfig.Nickname)

	// fmtJson, _ := json.MarshalIndent(CurrentConfig, "", "  ")
	// fmt.Println(string(fmtJson))

	sleepTime := time.Millisecond * 500
	sleepOffline := time.Millisecond * 1000
	sleepEventSwitch := time.Millisecond * 3000

	pipeFilePath := flag.String("pipe", "event.pipe", "Path to event.pipe")
	flag.Parse()

	fmt.Println("path:", *pipeFilePath)

	_, err = os.Stat(*pipeFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Invalid pipe file path", *pipeFilePath)
		return
	}

	pipeFile, err := os.OpenFile(*pipeFilePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening pipe:", err)
		return
	}

	defer pipeFile.Close()

	currentState := "idle"
	var newState string
	var player *client.Entity
	// var aircrafts *[]client.Entity
	var enemyAircrafts *[]client.Entity
	var enemyGroundUnits *[]client.Entity
	var nearestEnemyAir *client.Entity
	var nearestEnemyGround *client.Entity
	var distanceToEnemyAir float64
	var distanceToEnemyGround float64
	var lastDamage *client.Damage
	var hasCrashed, isShotDown, isMissionEnded bool
	var dangerBeginTime time.Time

	writeToFile(currentState, pipeFile)

	for {
		err = state.Load(CurrentConfig.Host)
		if err != nil {
			fmt.Println("state error: ", err, state)
			time.Sleep(sleepOffline)
			err = state.Load(CurrentConfig.Host)
			if err != nil {
				if currentState != "idle" {
					newState = "idle"
					writeToFile(newState, pipeFile)
					currentState = newState
				}
				continue
			}
			// } else {
			// fmt.Println(state)
		}
		err = indicators.Load(CurrentConfig.Host)
		if err != nil {
			fmt.Println("indicators error: ", indicators)
			// } else {
			// fmt.Println(indicators)
		}
		err = mapInfo.Load(CurrentConfig.Host)
		if err != nil {
			fmt.Println("mapInfo error: ", mapInfo)
			// } else {
			// fmt.Println(mapInfo)
		}

		if mapInfo.IsValid() {
			err = mapObj.Load(CurrentConfig.Host)
			if err != nil {
				fmt.Println("mapObj error: ", mapObj)
				// } else {
				// fmt.Println(mapObj)
			}
			err = hudMsg.Load(CurrentConfig.Host, lastEvt, lastDmg)
			if err != nil {
				fmt.Println("hudMsg error: ", mapObj)
			}

			player = mapObj.GetPlayerEntity()
			// aircrafts = mapObj.GetAircrafts()
			enemyAircrafts = mapObj.GetEnemyAircrafts(&CurrentConfig)
			enemyGroundUnits = mapObj.GetEnemyGroundUnits(&CurrentConfig)
			isShotDown = hudMsg.IsShotDown(CurrentConfig.Nickname)
			hasCrashed = hudMsg.HasCrashed(CurrentConfig.Nickname)
			isMissionEnded = hudMsg.IsMissionEnded()
			lastDamage = hudMsg.GetLastDmg()
			if lastDamage != nil {
				lastDmg = uint64(lastDamage.ID)
			}

			// fmt.Println("map loaded", currentState)
			switch currentState {
			case "start":
				// isShotDown = false
				// hasCrashed = false
				// isMissionEnded = false
				newState = "load"
				writeToFile(newState, pipeFile)
				currentState = newState
				time.Sleep(sleepEventSwitch)
			case "load", "idle":
				if mapObj.GetPlayerEntity() != nil {
					isShotDown = false
					hasCrashed = false
					isMissionEnded = false
					newState = "begin"
					writeToFile(newState, pipeFile)
					currentState = newState
					time.Sleep(sleepEventSwitch)
				}
			case "battle", "danger", "clear", "begin", "landed":
				// fmt.Println("player", player)
				if player == nil {
					if isShotDown || hasCrashed {
						newState = "death"
						writeToFile(newState, pipeFile)
						currentState = newState
						time.Sleep(sleepEventSwitch)
					} else if isMissionEnded {
						newState = "success"
						writeToFile(newState, pipeFile)
						currentState = newState
						time.Sleep(sleepEventSwitch)
					} else if currentState != "landed" {
						// probably landed
						newState = "landed"
						writeToFile(newState, pipeFile)
						currentState = newState
						time.Sleep(sleepEventSwitch)
					}
				} else {
					nearestEnemyAir = getNearestEntity(player, enemyAircrafts, &mapObj, &mapInfo)
					nearestEnemyGround = getNearestEntity(player, enemyGroundUnits, &mapObj, &mapInfo)
					if nearestEnemyAir != nil {
						distanceToEnemyAir = mapObj.GetDistance(player, nearestEnemyAir, &mapInfo)
						if distanceToEnemyAir <= 2000 {
							newState = "battle"
							dangerBeginTime = time.Now()
						} else if currentState != "battle" || distanceToEnemyAir > 10000 {
							newState = "danger"
							dangerBeginTime = time.Now()
						}
						if currentState != newState {
							writeToFile(newState, pipeFile)
							currentState = newState
							time.Sleep(sleepEventSwitch)
						}
					} else if nearestEnemyGround != nil {
						distanceToEnemyGround = mapObj.GetDistance(player, nearestEnemyGround, &mapInfo)
						if distanceToEnemyGround <= 20000 {
							newState = "danger"
							dangerBeginTime = time.Now()
						} else if time.Now().Unix()-dangerBeginTime.Unix() > 60 {
							newState = "clear"
						}
						if newState != currentState {
							writeToFile(newState, pipeFile)
							currentState = newState
							time.Sleep(sleepEventSwitch)
						}
					} else {
						if (currentState == "danger" || currentState == "battle") && time.Now().Unix()-dangerBeginTime.Unix() > 60 {
							newState = "clear"
							writeToFile(newState, pipeFile)
							currentState = newState
							time.Sleep(sleepEventSwitch)
						}
					}
				}
			case "death":
				if player != nil {
					newState = "begin"
					writeToFile(newState, pipeFile)
					currentState = newState
					time.Sleep(sleepEventSwitch)
				}
			}
			// fmt.Println("lastDmg", lastDmg, isShotDown, hasCrashed, isMissionEnded)
		} else {
			enemyAircrafts = nil
			isMissionEnded = true
			switch currentState {
			case "idle", "begin", "danger", "clear", "success", "death", "failure":
				newState = "start"
				writeToFile(newState, pipeFile)
				currentState = newState
				time.Sleep(sleepEventSwitch)
			}
		}

		time.Sleep(sleepTime)
	}
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
