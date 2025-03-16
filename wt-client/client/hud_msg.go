package client

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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

const playerShotDownRegExpTemplate = `%s.+(shot\s+down|destroyed).+` // [^\(]+\([^\)]+\)
// const playerShotDownRegExpTemplate = `%s\s+\([^\)]+\)\s+(shot\s+down|destroyed)\s+` // [^\(]+\([^\)]+\)

// const playerShootDownEnemyRegexp = `[^\(]+\s+\([^\)]+\)\s+shot\s+down\s+[^\(]+\([^\)]+\)`

func (hudMsg *HudMsg) Unmarshal(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &hudMsg)
}

func (hudMsg *HudMsg) Load(host string, lastEvt uint64, lastDmg uint64) error {
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

func (hudMsg *HudMsg) HasCrashed(nickname string) bool {
	for _, msg := range hudMsg.Damage {
		if strings.Contains(msg.Msg, "has crashed") && strings.Contains(msg.Msg, nickname) {
			return true
		}
	}

	return false
}

func (hudMsg *HudMsg) IsShotDown(nickname string) bool {
	for _, msg := range hudMsg.Damage {
		if strings.Contains(msg.Msg, "shot down "+nickname) {
			return true
		}
	}

	return false
}

func (hudMsg *HudMsg) HadKill(nickname string) bool {
	pattern := fmt.Sprintf(playerShotDownRegExpTemplate, nickname)
	for _, msg := range hudMsg.Damage {
		if matches, _ := regexp.MatchString(pattern, msg.Msg); matches {
			return true
		}
	}

	return false
}

func (hudMsg *HudMsg) GetLastKillTime(nickname string) int {
	var lastKillTime int
	pattern := fmt.Sprintf(playerShotDownRegExpTemplate, nickname)
	for _, msg := range hudMsg.Damage {
		if matches, _ := regexp.MatchString(pattern, msg.Msg); matches {
			lastKillTime = msg.Time
		}
	}

	return lastKillTime
}

func (hudMsg *HudMsg) IsMissionEnded() bool {
	for _, msg := range hudMsg.Damage {
		if strings.Contains(msg.Msg, "has achieved") || strings.Contains(msg.Msg, "has delivered the final blow!") {
			return true
		}
	}

	return false
}

func (hudMsg *HudMsg) GetLastDmg() *Damage {
	lastIndex := len(hudMsg.Damage) - 1
	if lastIndex < 0 {
		return nil
	}
	return &hudMsg.Damage[lastIndex]
}
