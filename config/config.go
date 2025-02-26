package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
)

type Config struct {
	Volume int              `json:"volume"`
	Path   string           `json:"path"`
	Events map[string]Event `json:"events"`
}

type Event struct {
	Volume int    `json:"volume"`
	Songs  []Song `json:"songs"`
}

type Song struct {
	Path      string `json:"path"`
	Skip      int    `json:"skip"`
	EndBefore int    `json:"endBefore"`
}

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Cannot open config file", path)
		return &Config{}, err
	}

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Cannot decode json", err)
		return &Config{}, err
	}

	return &config, nil
}

func (c Config) GetEvent(event string) (*Event, error) {
	ev, ok := c.Events[event]
	if ok {
		return &ev, nil
	} else {
		return &Event{}, errors.New("event not found")
	}
}

func (c Config) GetEventSong(event string, index uint8) (*Song, error) {
	ev, err := c.GetEvent(event)
	if err != nil {
		return &Song{}, err
	}

	song := ev.Songs[index]
	return &song, nil
}

func (c Config) GetRandomSong(event string) (*Song, error) {
	ev, err := c.GetEvent(event)
	if err != nil {
		return &Song{}, err
	}
	var index uint8
	length := len(ev.Songs)
	if length > 1 {
		index = uint8(rand.Intn(int(length - 1)))
	} else {
		index = 0
	}

	song, err := c.GetEventSong(event, index)

	return song, err
}

func (c Config) EventExists(event string) bool {
	_, ok := c.Events[event]

	return ok
}
