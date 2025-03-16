package model

import (
	"errors"
	"fmt"
	"math/rand"
)

type Event struct {
	Volume int   `json:"volume"`
	Sfx    []Sfx `json:"sfx"`
}

func (e *Event) GetRandomSfx() (*Sfx, error) {
	var index uint8
	length := len(e.Sfx)
	if length > 1 {
		index = uint8(rand.Intn(length))
		// ui.Debug("random sfx index:", index)
	} else {
		index = 0
	}
	return e.GetSfx(index)
}

func (e *Event) GetSfx(index uint8) (*Sfx, error) {
	if int(index) > len(e.Sfx) {
		return nil, fmt.Errorf("sfx index %d does not exists", index)
	}
	sfx := e.Sfx[index]
	if sfx.Volume == 0 {
		sfx.Volume = uint8(e.Volume)
	}

	return &sfx, nil
}

func (c *Config) GetEvent(event string) (*Event, error) {
	ev, ok := c.Events[event]
	if ok {
		return &ev, nil
	} else {
		return &Event{}, errors.New("event not found")
	}
}

func (c *Config) EventExists(event string) bool {
	_, ok := c.Events[event]

	return ok
}
