package model

import (
	"errors"
	"math/rand"
)

type Event struct {
	Volume *int  `json:"volume"`
	Sfx    []Sfx `json:"sfx"`
}

func (c *Config) GetEvent(event string) (*Event, error) {
	ev, ok := c.Events[event]
	if ok {
		return &ev, nil
	} else {
		return &Event{}, errors.New("event not found")
	}
}

func (c *Config) GetEventSfx(event string, index uint8) (*Sfx, error) {
	ev, err := c.GetEvent(event)
	if err != nil {
		return &Sfx{}, err
	}

	sfx := ev.Sfx[index]
	return &sfx, nil
}

func (c *Config) GetRandomEventSfx(event string) (*Sfx, error) {
	ev, err := c.GetEvent(event)
	if err != nil {
		return &Sfx{}, err
	}
	var index uint8
	length := len(ev.Sfx)
	if length > 1 {
		index = uint8(rand.Intn(length))
		// ui.Debug("random sfx index:", index)
	} else {
		index = 0
	}

	song, err := c.GetEventSfx(event, index)

	return song, err
}

func (c *Config) EventExists(event string) bool {
	_, ok := c.Events[event]

	return ok
}
