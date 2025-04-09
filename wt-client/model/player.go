package model

import (
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
)

type PlayerPosition struct {
	px   float64
	py   float64
	time int64
}

type PlayerSpeed struct {
	dx   float64
	dy   float64
	time int64
}

type ComputedSpeed struct {
	value float64
	time  int64
}

type Player struct {
	Name                    string
	Vehicle                 string
	Dead                    bool
	Damaged                 bool
	SeverlyDamaged          bool
	LastKillTime            int64
	LastDamageTime          int64
	LastSeverDamageTime     int64
	LastBurnTime            int64
	LastKilledTime          int64
	LastDamagedTime         int64
	LastSeverelyDamagedTime int64
	LastBurnedTime          int64
	Landed                  bool
	Targets                 map[string]*Player
	CurrentTarget           *Player
	CurrentEntity           *client.Entity
	LastPosition            *PlayerPosition
	LastComputedSpeed       *ComputedSpeed
	Positions               []*PlayerPosition
	IsDrone                 bool
	IsGround                bool
	IsAir                   bool
}

func (p *Player) LoadData(entity *client.Entity, nearestAirfield *client.Entity, state *client.State, indicators *client.Indicators, mapInfo *client.MapInfo, mapObj *client.MapObj) {
	p.CurrentEntity = entity
	currentPosition := PlayerPosition{entity.X, entity.Y, time.Now().UnixMilli()}
	p.Positions = append(p.Positions, p.LastPosition)

	var computedSpeed ComputedSpeed
	if p.LastPosition != nil {
		duration := time.Duration(currentPosition.time-p.LastPosition.time) * time.Millisecond
		// fmt.Println(duration)
		speedValue := mapInfo.GetSpeedKmh(p.LastPosition.px, p.LastPosition.py, currentPosition.px, currentPosition.py, duration)
		computedSpeed = ComputedSpeed{
			value: speedValue,
			time:  currentPosition.time,
		}
		// fmt.Println(speedValue)
	} else {
		computedSpeed = ComputedSpeed{
			value: 0,
			time:  currentPosition.time,
		}
	}
	p.LastComputedSpeed = &computedSpeed
	p.LastPosition = &currentPosition

	p.IsDrone = !indicators.Valid && state.IsValid()
	p.IsGround = indicators.Valid && indicators.Army == "tank"
	p.IsAir = indicators.Valid && indicators.Army == "air"
	if nearestAirfield == nil {
		p.Landed = false
	} else {
		p.checkLanded(nearestAirfield, mapInfo, mapObj)
	}

	fmt.Println("isGround", p.IsGround)
	fmt.Println("isAir", p.IsAir)
	fmt.Println("isDrone", p.IsDrone)
	fmt.Println("Landed", p.Landed)
}

func (p *Player) checkLanded(airfield *client.Entity, mapInfo *client.MapInfo, mapObj *client.MapObj) {
	airfield.X = airfield.Sx
	airfield.Y = airfield.Sy
	fmt.Println("afd: ", mapObj.GetDistance(p.CurrentEntity, airfield, mapInfo))
	p.Landed = (p.LastComputedSpeed.value == 0 || !p.LastComputedSpeed.IsValid()) && p.IsAir && mapObj.GetDistance(p.CurrentEntity, airfield, mapInfo) < 2000
}

func (cs *ComputedSpeed) IsValid() bool {
	return cs.value < 2000
}
