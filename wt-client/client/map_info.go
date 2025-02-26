package client

import (
	"encoding/json"
)

type MapInfo struct {
	Valid         bool      `json:"valid"`
	GridSize      []float64 `json:"grid_size"`
	GridSteps     []float64 `json:"grid_steps"`
	GridZero      []float64 `json:"grid_zero"`
	HudType       int       `json:"hud_type"`
	MapGeneration int       `json:"map_generation"`
	MapMax        []float64 `json:"map_max"`
	MapMin        []float64 `json:"map_min"`
}

func (mapInfo *MapInfo) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &mapInfo)
}

func (mapInfo *MapInfo) IsValid() bool {
	return mapInfo.Valid
}

func (mapInfo *MapInfo) Load(host string) error {
	body, err := GetDataFromUrl(host + "map_info.json")
	if err != nil {
		return err
	}

	err = mapInfo.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}
