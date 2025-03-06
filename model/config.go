package model

import (
	"encoding/json"
	"os"
)

type Config struct {
	PlayerType string           `json:"player_type"`
	ServerType string           `json:"server_type"`
	Volume     uint8            `json:"volume"`
	Path       string           `json:"path"`
	Events     map[string]Event `json:"events"`
	States     map[string]State `json:"states"`
}

func (c *Config) Read(path string) error {
	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
