package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Nickname string `json:"nickname"`
	Colors   struct {
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
