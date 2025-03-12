package client

import (
	"regexp"
	"testing"
)

func TestRegexes(t *testing.T) {
	pattern := `[^\(]+\s+\([^\)]+\)\s+shot\s+down\s+[^\(]+\([^\)]+\)`
	text1 := `-Tygri- xbladecz (F-20A) shot down Meinhard (MiG-23M)`
	// text2 := `Vsevolod (F-5E) shot down -Tygri- xbladecz (F-20A)`

	matched, err := regexp.MatchString(pattern, text1)

	if err != nil {
		t.Error(err)
	}
	if !matched {
		t.Errorf("pattern is wrong %s", pattern)
	}

	// matched, _ = regexp.MatchString(pattern, text2)

	// if matched {
	// 	t.Errorf("pattern is wrong %s", pattern)
	// }

}
