package client

import (
	"encoding/json"
	"math"
	"time"
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

func (m *MapInfo) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &m)
}

func (m *MapInfo) IsValid() bool {
	return m.Valid
}

func (m *MapInfo) Load(host string) error {
	body, err := GetDataFromUrl(host + "map_info.json")
	if err != nil {
		return err
	}

	err = m.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}

func (m *MapInfo) GetDistance(x1, y1, x2, y2 float64) float64 {
	dx := (x2 - x1) * (m.MapMax[0] - m.MapMin[0])
	dy := (y2 - y1) * (m.MapMax[1] - m.MapMin[1])

	return math.Sqrt(dx*dx + dy*dy)
}

func (m *MapInfo) GetSpeed(x1, y1, x2, y2 float64, duration time.Duration) float64 {
	d := m.GetDistance(x1, y1, x2, y2)
	return d / duration.Seconds()
}

func (m *MapInfo) GetSpeedKmh(x1, y1, x2, y2 float64, duration time.Duration) float64 {
	speed := m.GetSpeed(x1, y1, x2, y2, duration)

	return speed * 3.6
}
