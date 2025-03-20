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

func (hudMsg *HudMsg) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &hudMsg)
}

func (hudMsg *HudMsg) Load(host string, lastEvt uint64, lastDmg uint64) error {
	// fmt.Printf("%vhudmsg?lastEvt=%v&lastDmg=%v", host, lastEvt, lastDmg)
	body, err := GetDataFromUrl(fmt.Sprintf("%vhudmsg?lastEvt=%v&lastDmg=%v", host, lastEvt, lastDmg))
	if err != nil {
		return err
	}

	err = hudMsg.Unmarshal(body)
	if err != nil {
		return err
	}

	return nil
}

func (hudMsg *HudMsg) MatchMessages(pattern *regexp.Regexp) []Damage {
	// fmt.Println("c: ", len(hudMsg.Damage), pattern.String())
	matched := make([]Damage, 0)
	for _, msg := range hudMsg.Damage {
		if pattern.MatchString(msg.Msg) {
			matched = append(matched, msg)
		}
	}

	return matched
}

func (hudMsg *HudMsg) GetLastDmg() *Damage {
	lastIndex := len(hudMsg.Damage) - 1
	if lastIndex < 0 {
		return nil
	}
	return &hudMsg.Damage[lastIndex]
}
