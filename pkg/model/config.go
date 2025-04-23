package model

import (
	"encoding/json"
	"os"

	"github.com/lamasutra/bg-music/pkg/input"
)

type InputDevice struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (i *InputDevice) ToManagableDevice() input.InputDevice {
	return input.InputDevice{Name: i.Name, Path: i.Path}
}

type InputDeviceControls struct {
	Device  InputDevice       `json:"device"`
	Mapping map[string]string `json:"mapping"`
}

type Config struct {
	PlayerType string              `json:"player_type"`
	ServerType string              `json:"server_type"`
	Controls   InputDeviceControls `json:"controls"`
	Volume     uint8               `json:"volume"`
	Path       string              `json:"path"`
	Events     map[string]Event    `json:"events"`
	States     map[string]State    `json:"states"`
	Narrate    map[string]Speech   `json:"narrate"`
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

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0755)
}

func (c *Config) GetSfxPath(sfx *Sfx) string {
	return c.Path + "/" + sfx.Path
}

func (c *Config) GetMusicPath(music *Music) string {
	return c.
		Path + "/" + music.Path
}

func (c *Config) GetInputDeviceControls() InputDeviceControls {
	return c.Controls
}

func (c *Config) SetInputDeviceControls(controls InputDeviceControls) {
	c.Controls = controls
}
