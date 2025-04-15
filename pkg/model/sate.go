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

func (st *State) GetMusic(index uint8) (*Music, error) {
	music := st.Music[index]
	if music.Volume == 0 {
		music.Volume = *st.Volume
	}
	return &music, nil
}

func (c *Config) GetStateMusic(state string, index uint8) (*Music, error) {
	st, err := c.GetState(state)
	if err != nil {
		return &Music{}, err
	}

	return st.GetMusic(index)
}

func (c *Config) GetStatePlaylist(state string) (*[]Music, error) {
	st, err := c.GetState(state)
	if err != nil {
		return nil, nil
	}

	return &st.Music, nil
}

func (st *State) GetRandomMusic() (*Music, error) {
	var index uint8
	length := len(st.Music)
	if length > 1 {
		index = uint8(rand.Intn(length))
	} else {
		index = 0
	}

	music, err := st.GetMusic(index)

	return music, err
}

func (c *Config) GetRandomStateMusic(state string) (*Music, error) {
	st, err := c.GetState(state)
	if err != nil {
		return &Music{}, err
	}
	return st.GetRandomMusic()
}

func (c *Config) StateExists(state string) bool {
	_, ok := c.States[state]

	return ok
}
