package input

import (
	"fmt"
	"regexp"
	"testing"
)

func TestPatterns(t *testing.T) {
	nickname := "-Tygri- xbladecz"
	strCrash := "-Tygri- xbladecz (F/A-18C Late) has crashed."
	strShotDown := "=FLXG= 3472254422 (Su-30SM) shot down -Tygri- xbladecz (F/A-18C Late)"
	strAnyKill := "=CPNB2= J35_Goshawk鶻鹰 (F-16C) shot down -NFlyZ- End and Death (J-11B)"
	strKill := "-Tygri- xbladecz (F/A-18C Late) shot down =CPNB2= J35_Goshawk鶻鹰 (F-16C)"
	playerMadeKilledRegExp = *regexp.MustCompile(fmt.Sprintf(playerMadeKillPatternTemplate, nickname))
	anyKillRegExp = *regexp.MustCompile(anyKillPattern)
	playerHasCrashedRegExp = *regexp.MustCompile(fmt.Sprintf(playerHasCrashedPatternTemplate, nickname))
	playerIsShotDownRegExp = *regexp.MustCompile(fmt.Sprintf(playerIsShotDownPatternTemplate, nickname))

	if !playerMadeKilledRegExp.MatchString(strKill) {
		t.Error("player kill pattern is invalid")
	}
	if playerMadeKilledRegExp.MatchString(strAnyKill) {
		t.Error("player kill pattern is invalid")
	}
	if !anyKillRegExp.MatchString(strAnyKill) {
		t.Error("any kill pattern is invalid")
	}
	if !anyKillRegExp.MatchString(strKill) {
		t.Error("any kill pattern is invalid")
	}

	if !playerHasCrashedRegExp.MatchString(strCrash) {
		t.Error("player has crashed pattern is invalid")
	}
	if playerHasCrashedRegExp.MatchString(strShotDown) {
		t.Error("player has crashed pattern is invalid")
	}
	if !playerIsShotDownRegExp.MatchString(strShotDown) {
		t.Error("player is shot down pattern is invalid")
	}
	if playerIsShotDownRegExp.MatchString(strCrash) {
		t.Error("player is shot down pattern is invalid")
	}
}
