package clientConfig

import (
	"encoding/json"
	"fmt"
	"os"
)

type Sfx struct {
	Path string `json:"path"`
}

type Event struct {
	Volume *int  `json:"volume"`
	Sfx    []Sfx `json:"sfx"`
}

type Vehicle struct {
	Title string `json:"title"`
	Theme string `json:"theme"`
}

type Map struct {
	Title    string `json:"title"`
	Checksum uint32 `json:"crc32"`
}

type EventStates struct {
	Events map[string]Event `json:"events"`
	States map[string]State `json:"states"`
}

type Theme struct {
	Title string
	EventStates
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
