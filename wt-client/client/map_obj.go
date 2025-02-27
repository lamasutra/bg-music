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
	dx := (ent2.X - ent1.X) * (mapInfo.MapMax[0] - mapInfo.MapMin[0])
	dy := (ent2.Y - ent1.Y) * (mapInfo.MapMax[1] - mapInfo.MapMin[1])

	// fmt.Println(dx, dy)

	return math.Sqrt(dx*dx + dy*dy)
}

func (mapObj *MapObj) GetAircrafts() *[]Entity {
	return mapObj.getEntitiesByType("aircraft")
}

func (mapObj *MapObj) GetEnemyAircrafts(colors *[]string) *[]Entity {
	return mapObj.getEntitiesByTypeAndColors("aircraft", colors)
}

func (mapObj *MapObj) GetEnemyGroundUnits(colors *[]string) *[]Entity {
	return mapObj.getEntitiesByTypeAndColors("ground_model", colors)
}

func (mapObj *MapObj) GetTanks() *[]Entity {
	return mapObj.getEntitiesByType("ground_model")
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
