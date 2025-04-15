package model

import (
	"fmt"
	"strings"
)

type Vehicle struct {
	Title  string `json:"title"`
	Theme  string `json:"theme"`
	Type   string `json:"type"`
	Volume int    `json:"volume"`
}

var nameTypeCache = make(map[string]string, 24)

type VehicleList []Vehicle

func (v *Vehicle) merge(from Vehicle) *Vehicle {
	merged := Vehicle{}
	if from.Theme != "" {
		merged.Theme = from.Theme
	} else {
		merged.Theme = v.Theme
	}
	if from.Title != "" {
		merged.Title = from.Title
	} else {
		merged.Title = v.Title
	}
	if from.Volume != 0 {
		merged.Volume = from.Volume
	} else {
		merged.Volume = v.Volume
	}
	if from.Type != "" {
		merged.Type = from.Type
	} else {
		merged.Type = v.Type
	}

	return &merged
}

func (vl *VehicleList) DetectType(name string) (string, error) {
	vType, ok := nameTypeCache[name]
	if ok {
		return vType, nil
	}

	for _, v := range *vl {
		if strings.Contains(name, v.Title) {
			nameTypeCache[name] = v.Type
			return v.Type, nil
		}
	}

	return "unknown", fmt.Errorf("unknonwn type `%s`", name)
}
