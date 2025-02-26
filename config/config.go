package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/bg-music/model"
)

type Config struct {
	Volume uint8                  `json:"volume"`
	Path   string                 `json:"path"`
	Events map[string]model.Event `json:"events"`
	States map[string]model.State `json:"states"`
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

func (c Config) GetEvent(event string) (*model.Event, error) {
	ev, ok := c.Events[event]
	if ok {
		return &ev, nil
	} else {
		return &model.Event{}, errors.New("event not found")
	}
}

func (c Config) GetEventSong(event string, index uint8) (*model.Music, error) {
	ev, err := c.GetEvent(event)
	if err != nil {
		return &model.Music{}, err
	}

	song := ev.Music[index]
	return &song, nil
}

func (c Config) GetRandomSong(event string) (*model.Music, error) {
	ev, err := c.GetEvent(event)
	if err != nil {
		return &model.Music{}, err
	}
	var index uint8
	length := len(ev.Music)
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
