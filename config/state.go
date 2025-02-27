package config

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/lamasutra/bg-music/model"
)

func (c *Config) GetState(state string) (*model.State, error) {
	st, ok := c.States[state]
	if ok {
		return &st, nil
	} else {
		return &model.State{}, errors.New("state not found")
	}
}

func (c *Config) GetStateMusic(state string, index uint8) (*model.Music, error) {
	ev, err := c.GetState(state)
	if err != nil {
		return &model.Music{}, err
	}

	music := ev.Music[index]
	return &music, nil
}

func (c *Config) GetRandomStateMusic(state string) (*model.Music, error) {
	st, err := c.GetState(state)
	if err != nil {
		return &model.Music{}, err
	}
	var index uint8
	length := len(st.Music)
	if length > 1 {
		index = uint8(rand.Intn(length))
		fmt.Println("rand:", index)
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
