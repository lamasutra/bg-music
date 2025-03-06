package model

import (
	"errors"
	"math/rand"
)

type State struct {
	Volume *int     `json:"volume"`
	Music  []Music  `json:"music"`
	States []string `json:"states"`
}

func (c *Config) GetState(state string) (*State, error) {
	st, ok := c.States[state]
	if ok {
		return &st, nil
	} else {
		return &State{}, errors.New("state not found")
	}
}

func (c *Config) GetStateMusic(state string, index uint8) (*Music, error) {
	ev, err := c.GetState(state)
	if err != nil {
		return &Music{}, err
	}

	music := ev.Music[index]
	return &music, nil
}

func (c *Config) GetRandomStateMusic(state string) (*Music, error) {
	st, err := c.GetState(state)
	if err != nil {
		return &Music{}, err
	}
	var index uint8
	length := len(st.Music)
	if length > 1 {
		index = uint8(rand.Intn(length))
		// ui.Debug("random music index:", index)
	} else {
		index = 0
	}

	music, err := c.GetStateMusic(state, index)

	return music, err
}

func (c *Config) StateExists(state string) bool {
	_, ok := c.States[state]

	return ok
}
