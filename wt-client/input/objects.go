package input

type location struct {
	Id int64
	X  float64
	Y  float64
	Dx float64
	Dy float64
}

type vehicle struct {
	location
	Type  string
	Speed float64
}

type Aircraft struct {
	vehicle
}

type Airfield struct {
	location
}

type GroundVehicle struct {
	vehicle
}

type CaptureZone struct {
	location
	State string
}

type BombingPoint struct {
	location
	State string
}

type DefendPoint struct {
	location
}

type MapObjects struct {
	DefendPoints   []DefendPoint
	BombingPoints  []BombingPoint
	CaptureZones   []CaptureZone
	GroundVehicles struct {
		Friendly []GroundVehicle
		Foe      []GroundVehicle
	}
	AirVehicles struct {
		Friendly []GroundVehicle
		Foe      []GroundVehicle
	}
}
