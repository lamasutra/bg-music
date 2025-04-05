package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lamasutra/bg-music/wt-client/ui"
)

var mergedThemesCache map[string]*Theme = make(map[string]*Theme)

type Sfx struct {
	Path        string `json:"path"`
	SampleRate  int64  `json:"sample_rate"`
	NumChannels int    `json:"num_channels"`
	Volume      uint8  `json:"volume"`
}

type Map struct {
	Title    string `json:"title"`
	Checksum uint32 `json:"crc32"`
	Theme    string `json:"theme"`
}

type Music struct {
	Path      string `json:"path"`
	Skip      int    `json:"skip"`
	EndBefore int    `json:"endBefore"`
}

type Config struct {
	Nickname     string             `json:"nickname"`
	Host         string             `json:"host"`
	BgPlayerType string             `json:"bg_player_type"`
	BgPlayerHost string             `json:"bg_player_host"`
	Themes       map[string]Theme   `json:"themes"`
	Maps         map[string]Map     `json:"maps"`
	Vehicles     map[string]Vehicle `json:"vehicles"`
	StateRules   StateRules         `json:"state_rules"`
	Colors       struct {
		Friendly struct {
			Air    []string `json:"air"`
			Ground []string `json:"ground"`
		} `json:"friendly"`
		Foe struct {
			Air    []string `json:"air"`
			Ground []string `json:"ground"`
		} `json:"foe"`
	} `json:"colors"`
}

func (config *Config) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &config)
}

func (config *Config) Read(path string) error {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Cannot open config file", path)
		return err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Cannot decode json", err)
		return err
	}

	return nil
}

func (c *Config) getVehicleConfig(vehicle string) *Vehicle {
	conf, ok := c.Vehicles[vehicle]
	if ok {
		ui.Debug("config for vehicle `", vehicle, "` not found")
		return &conf
	}
	ui.Debug("vehicle config for `", vehicle, "` not found")
	return &Vehicle{}
}

func (c *Config) getTheme(theme string) *Theme {
	ui.Debug("getting theme", theme)
	found, exists := c.Themes[theme]
	if !exists {
		ui.Error("vehicle theme", theme, "does not exists")
		return &Theme{}
	}
	// fmt.Println("reference", fmt.Sprintf("%p", &found))
	if found.Extend != "" {
		extend, exists := c.Themes[found.Extend]
		ui.Debug("extending", found.Extend)
		// fmt.Println("reference extend", fmt.Sprintf("%p", &extend))
		if exists {
			found = *extend.merge(found)
			ui.Debug("extended")
			// fmt.Println("reference merge", fmt.Sprintf("%p", &found))
		} else {
			ui.Error("does not exists", found.Extend)
		}
	} else {
		ui.Debug("theme is not extended")
	}

	return &found
}

func (c *Config) GetVehicleForPlayerTypeAndVehicleType(playerType string, vehicleType string) *Vehicle {
	typeConf := c.getVehicleConfig(playerType)
	vehicleConf := c.getVehicleConfig(vehicleType)
	if vehicleConf.Theme == "" {
		vehicleConf.Theme = "default"
	}

	conf := vehicleConf.merge(*typeConf)

	if conf.Title == "" {
		conf.Title = vehicleType
	}

	return conf
}

func (c *Config) GetThemeForVehicle(vehicle *Vehicle) *Theme {
	ui.Debug("GetThemeForVehicle", vehicle)
	cacheKey := vehicle.Title
	if cacheKey == "" {
		cacheKey = vehicle.Type
	}
	// fmt.Println("cache key", cacheKey)
	theme, ok := mergedThemesCache[cacheKey]
	if ok {
		// js, _ := json.MarshalIndent(mergedThemesCache[cacheKey], "", "  ")
		// fmt.Println("found in cache", string(js))
		return theme
	}

	theme = c.getTheme(vehicle.Theme)
	vehicleTheme := *theme
	// fmt.Println("references", fmt.Sprintf("%p, %p", theme, &vehicleTheme))

	if vehicle.Volume > 0 {
		vehicleTheme.States = vehicleTheme.forceStateVolume(vehicle.Volume)
	}

	mergedThemesCache[cacheKey] = &vehicleTheme
	// fmt.Println("storing")

	return &vehicleTheme
}
