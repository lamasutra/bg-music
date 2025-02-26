package client

import (
	"encoding/json"
)

type Indicators struct {
	Valid            bool    `json:"valid"`
	Army             string  `json:"army"`
	Type             string  `json:"type"`
	Speed            float64 `json:"speed"`
	Pedals1          float64 `json:"pedals1"`
	Pedals2          float64 `json:"pedals2"`
	Pedals3          float64 `json:"pedals3"`
	Pedals4          float64 `json:"pedals4"`
	StickElevator    float64 `json:"stick_elevator"`
	StickAilerons    float64 `json:"stick_ailerons"`
	Vario            float64 `json:"vario"`
	AltitudeHour     float64 `json:"altitude_hour"`
	AltitudeMin      float64 `json:"altitude_min"`
	Altitude10K      float64 `json:"altitude_10k"`
	Altitude1Hour    float64 `json:"altitude1_hour"`
	Altitude1Min     float64 `json:"altitude1_min"`
	Altitude1100     float64 `json:"altitude1_10k"`
	AviahorizonRoll  float64 `json:"aviahorizon_roll"`
	AviahorizonPitch float64 `json:"aviahorizon_pitch"`
	Bank             float64 `json:"bank"`
	Bank2            float64 `json:"bank2"`
	Turn             float64 `json:"turn"`
	Compass1         float64 `json:"compass1"`
	Compass2         float64 `json:"compass2"`
	ClockHour        float64 `json:"clock_hour"`
	ClockMin         float64 `json:"clock_min"`
	ClockSec         float64 `json:"clock_sec"`
	ManifoldPressure float64 `json:"manifold_pressure"`
	RpmMin           float64 `json:"rpm_min"`
	RpmHour          float64 `json:"rpm_hour"`
	OilPressure      float64 `json:"oil_pressure"`
	OilPressure1     float64 `json:"oil_pressure1"`
	OilTemperature   float64 `json:"oil_temperature"`
	HeadTemperature  float64 `json:"head_temperature"`
	Mixture          float64 `json:"mixture"`
	CarbTemperature  float64 `json:"carb_temperature"`
	Fuel1            float64 `json:"fuel1"`
	FuelPressure     float64 `json:"fuel_pressure"`
	Gears            float64 `json:"gears"`
	Gears1           float64 `json:"gears1"`
	Flaps            float64 `json:"flaps"`
	Throttle         float64 `json:"throttle"`
	Weapon2          float64 `json:"weapon2"`
	Weapon3          float64 `json:"weapon3"`
	Weapon4          float64 `json:"weapon4"`
	PropPitch        float64 `json:"prop_pitch"`
	Supercharger     float64 `json:"supercharger"`
	FlapsIndicator   float64 `json:"flaps_indicator"`
	GearLIndicator   float64 `json:"gear_l_indicator"`
	GearRIndicator   float64 `json:"gear_r_indicator"`
	GearCIndicator   float64 `json:"gear_c_indicator"`
	GMeter           float64 `json:"g_meter"`
	GMeterMin        float64 `json:"g_meter_min"`
	GMeterMax        float64 `json:"g_meter_max"`
	Blister1         float64 `json:"blister1"`
	Blister2         float64 `json:"blister2"`
}

func (indicators *Indicators) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &indicators)
}

func (indicators *Indicators) IsValid() bool {
	return indicators.Valid
}

func (s *Indicators) Load(host string) error {
	body, err := GetDataFromUrl(host + "indicators")
	if err != nil {
		return err
	}

	err = s.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}
