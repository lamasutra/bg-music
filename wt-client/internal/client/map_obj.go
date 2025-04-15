package client

import (
	"encoding/json"
	"math"
	"strings"
)

type MapObj []Entity

type Entity struct {
	Type     string  `json:"type"`
	Color    string  `json:"color"`
	ColorRGB []int   `json:"color[]"`
	Blink    int     `json:"blink"`
	Icon     string  `json:"icon"`
	IconBG   string  `json:"icon_bg"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Dx       float64 `json:"dx"`
	Dy       float64 `json:"dy"`
	Sx       float64 `json:"sx"`
	Sy       float64 `json:"sy"`
	Ex       float64 `json:"ex"`
	Ey       float64 `json:"ey"`
}

func inArray(value string, haystack *[]string) bool {
	for _, val := range *haystack {
		if val == value {
			return true
		}
	}
	return false
}

func (mapObj *MapObj) getEntitiesByType(entType string) *[]Entity {
	var entities []Entity
	for _, entity := range *mapObj {
		if entity.Type == entType && entity.Icon != "Player" {
			entities = append(entities, entity)
		}
	}
	return &entities
}

func (mapObj *MapObj) getEntitiesByTypeAndColors(entType string, colors *[]string) *[]Entity {
	var entities []Entity
	for _, entity := range *mapObj {
		if entity.Type == entType && entity.Icon != "Player" && inArray(strings.ToLower(entity.Color), colors) {
			entities = append(entities, entity)
		}
	}
	return &entities
}

func (mapObj *MapObj) Unmarshal(jsonBytes []byte) error {
	*mapObj = nil
	return json.Unmarshal(jsonBytes, &mapObj)
}

func (mapObj *MapObj) GetPlayerEntity() *Entity {
	for _, entity := range *mapObj {
		if entity.Icon == "Player" {
			return &entity
		}
	}
	return nil
}

func (mapObj *MapObj) GetDistance(ent1 *Entity, ent2 *Entity, mapInfo *MapInfo) float64 {
	return mapInfo.GetDistance(ent1.X, ent1.Y, ent2.X, ent2.Y)
}

func (mapObj *MapObj) GetHeading(ent1 *Entity, ent2 *Entity) float64 {
	dx := ent2.X - ent1.X
	dy := ent2.Y - ent1.Y

	// Compute the angle in radians using atan2
	angleRad := math.Atan2(dy, dx)

	// Convert the angle from radians to degrees
	angleDeg := angleRad * 180 / math.Pi
	angleDeg += 90

	// Normalize the angle to be between 0° and 360°
	if angleDeg < 0 {
		angleDeg += 360
	}

	return angleDeg
}

func (mapObj *MapObj) GetAircrafts() *[]Entity {
	return mapObj.getEntitiesByType("aircraft")
}

func (mapObj *MapObj) GetAircraftsByColors(colors *[]string) *[]Entity {
	return mapObj.getEntitiesByTypeAndColors("aircraft", colors)
}

func (mapObj *MapObj) GetGroundUnitsByColors(colors *[]string) *[]Entity {
	return mapObj.getEntitiesByTypeAndColors("ground_model", colors)
}

func (mapObj *MapObj) GetTanks() *[]Entity {
	return mapObj.getEntitiesByType("ground_model")
}

func (mapObj *MapObj) GetTankRespawnBases() *[]Entity {
	return mapObj.getEntitiesByType("respawn_base_tank")
}

func (mapObj *MapObj) GetFighterRespawnBases() *[]Entity {
	return mapObj.getEntitiesByType("respawn_base_fighter")
}

func (mapObj *MapObj) GetCaptureZones() *[]Entity {
	return mapObj.getEntitiesByType("capture_zone")
}

func (mapObj *MapObj) GetAirfields() *[]Entity {
	return mapObj.getEntitiesByType("airfield")
}

func (mapObj *MapObj) Load(host string) error {
	body, err := GetDataFromUrl(host + "map_obj.json")
	if err != nil {
		return err
	}

	err = mapObj.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}
