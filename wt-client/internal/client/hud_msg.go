package client

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type HudMsg struct {
	Events []Event  `json:"events"`
	Damage []Damage `json:"damage"`
}

type Event struct{}

type Damage struct {
	ID     int    `json:"id"`
	Msg    string `json:"msg"`
	Sender string `json:"sender"`
	Enemy  bool   `json:"enemy"`
	Mode   string `json:"mode"`
	Time   int    `json:"time"`
}

func (h *HudMsg) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &h)
}

func (h *HudMsg) Load(host string, lastEvt uint64, lastDmg uint64) error {
	// fmt.Printf("%vhudmsg?lastEvt=%v&lastDmg=%v", host, lastEvt, lastDmg)
	body, err := GetDataFromUrl(fmt.Sprintf("%vhudmsg?lastEvt=%v&lastDmg=%v", host, lastEvt, lastDmg))
	if err != nil {
		return err
	}

	err = h.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}

func (h *HudMsg) Each(callback func(dmg Damage, index int) bool) {
	var brk bool
	for i, dmg := range h.Damage {
		brk = callback(dmg, i)
		if brk {
			break
		}
	}
}

func (h *HudMsg) MatchMessages(pattern *regexp.Regexp) []Damage {
	// fmt.Println("c: ", len(hudMsg.Damage), pattern.String())
	matched := make([]Damage, 0)
	for _, msg := range h.Damage {
		if pattern.MatchString(msg.Msg) {
			matched = append(matched, msg)
		}
	}

	return matched
}

func (h *HudMsg) GetLastDmg() *Damage {
	lastIndex := len(h.Damage) - 1
	if lastIndex < 0 {
		return nil
	}
	return &h.Damage[lastIndex]
}
