package client

import (
	"encoding/json"
)

type State struct {
	Valid             bool    `json:"valid"`
	Aileron           float64 `json:"aileron, %"`
	Elevator          float64 `json:"elevator, %"`
	Rudder            float64 `json:"rudder, %"`
	Flaps             float64 `json:"flaps, %"`
	Gear              float64 `json:"gear, %"`
	Airbrake          float64 `json:"airbrake, %"`
	H                 float64 `json:"H, m"`
	TAS               float64 `json:"TAS, km/h"`
	IAS               float64 `json:"IAS, km/h"`
	M                 float64 `json:"M"`
	AoA               float64 `json:"AoA, deg"`
	AoS               float64 `json:"AoS, deg"`
	Ny                float64 `json:"Ny"`
	Vy                float64 `json:"Vy, m/s"`
	Wx                float64 `json:"Wx, deg/s"`
	Mfuel             float64 `json:"Mfuel, kg"`
	Mfuel0            float64 `json:"Mfuel0, kg"`
	Throttle1         float64 `json:"throttle 1, %"`
	RPMThrottle1      float64 `json:"RPM throttle 1, %"`
	Mixture1          float64 `json:"mixture 1, %"`
	Radiator1         float64 `json:"radiator 1, %"`
	CompressorStage1  int     `json:"compressor stage 1"`
	Magneto1          int     `json:"magneto 1"`
	Power1            float64 `json:"power 1, hp"`
	RPM1              float64 `json:"RPM 1"`
	ManifoldPressure1 float64 `json:"manifold pressure 1, atm"`
	OilTemp1          float64 `json:"oil temp 1, C"`
	Pitch1            float64 `json:"pitch 1, deg"`
	Thrust1           float64 `json:"thrust 1, kgs"`
	Efficiency1       float64 `json:"efficiency 1, %"`
}

func (state *State) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &state)
}

func (state *State) IsValid() bool {
	return state.Valid
}

func (state *State) Load(host string) error {
	body, err := GetDataFromUrl(host + "state")
	if err != nil {
		return err
	}

	err = state.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}
