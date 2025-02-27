package config

import (
	"encoding/json"
	"os"

	"github.com/lamasutra/bg-music/model"
)

type Config struct {
	PlayerType string                 `json:"player_type"`
	ServerType string                 `json:"server_type"`
	Volume     uint8                  `json:"volume"`
	Path       string                 `json:"path"`
	Events     map[string]model.Event `json:"events"`
	States     map[string]model.State `json:"states"`
}

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return &Config{}, err
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}
