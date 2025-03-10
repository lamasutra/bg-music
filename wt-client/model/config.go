package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type Sfx struct {
	Path string `json:"path"`
}

type Vehicle struct {
	Title  string         `json:"title"`
	Theme  string         `json:"theme"`
	Volume map[string]int `json:"volume"`
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
	Nickname     string               `json:"nickname"`
	Host         string               `json:"host"`
	BgPlayerType string               `json:"bg_player_type"`
	Themes       map[string]Theme     `json:"themes"`
	Maps         map[string]Map       `json:"maps"`
	Vehicles     map[string]Vehicle   `json:"vehicles"`
	StateRules   map[string]StateRule `json:"state_rules"`
	Colors       struct {
		Enemy struct {
			Air    []string `json:"air"`
			Ground []string `json:"ground"`
		} `json:"enemy"`
		Friendly struct {
			Air    []string `json:"air"`
			Ground []string `json:"ground"`
		} `json:"friendly"`
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
		return &conf
	}
	return &Vehicle{}
}

func (c *Config) getTheme(theme string) *Theme {
	found, exists := c.Themes[theme]
	if !exists {
		return &Theme{}
	}

	return &found
}

func (c *Config) mergeThemes(from string, to string) *Theme {
	themeFrom := c.getTheme(from)
	themeTo := c.getTheme(to)

	return themeFrom.Merge(*themeTo)
}

func (c *Config) GetConfigForPlayerVehicle(playerType string, vehicle string) {
	typeConf := c.getVehicleConfig(playerType)
	vehicleConf := c.getVehicleConfig(vehicle)

	// typeConf.
}
