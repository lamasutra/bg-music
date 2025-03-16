package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"
)

const nickname = "-Tygri- xbladecz"

func TestPlayerShotDownRegExp(t *testing.T) {
	// text1 := `-Tygri- xbladecz (F-20A) shot down Someone (MiG-23M)`
	text1 := `-Tygri- xbladecz (F-14A) shot down -LEXIL- jannis19 (Su-27SM)`
	// text2 := `Vsevolod (F-5E) shot down -Tygri- xbladecz (F-20A)`

	pattern := fmt.Sprintf(playerShotDownRegExpTemplate, nickname)
	matched, err := regexp.MatchString(pattern, text1)

	if err != nil {
		t.Error(err)
	}
	if !matched {
		t.Errorf("pattern is wrong `%s`", pattern)
	}

	// matched, _ = regexp.MatchString(pattern, text2)

	// if matched {
	// 	t.Errorf("pattern is wrong %s", pattern)
	// }

}

func TestJson(t *testing.T) {
	hudmsg := HudMsg{}
	file, err := os.Open("test/hudmsg.json")
	if err != nil {
		t.Error(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(data, &hudmsg)

	killTime := hudmsg.GetLastKillTime(nickname)
	if killTime == 0 {
		t.Error("kill time not found")
	}
	// t.Error("kill time", killTime)
}
