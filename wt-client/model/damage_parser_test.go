package model

import (
	"testing"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

func TestDamageParser(t *testing.T) {
	ui.CreateUI("cli")
	d := DamageParser{}
	d.init()
	dmg := client.Damage{Msg: "-Tygri- xbladecz (F/A-18A) has crashed."}
	matched := d.parseDamage(dmg, 0)

	if matched == false {
		t.Error("did not match")
	}

	dmg = client.Damage{Msg: "Roland 1 shot down -Tygri- xbladecz (F/A-18A)"}
	matched = d.parseDamage(dmg, 0)

	if matched == false {
		t.Error("did not match")
	}

	dmg = client.Damage{Msg: "-Tygri- xbladecz (F-14A) critically damaged MiG-15bis ISh"}
	matched = d.parseDamage(dmg, 0)

	if matched == false {
		t.Error("did not match")
	}

	dmg = client.Damage{Msg: "-Tygri- xbladecz (F-14A) severely damaged MiG-15bis ISh"}
	matched = d.parseDamage(dmg, 0)

	if matched == false {
		t.Error("did not match")
	}

	dmg = client.Damage{Msg: "-Tygri- xbladecz (F-14A) destroyed [ai] MiG-15bis ISh"}
	matched = d.parseDamage(dmg, 0)

	if matched == false {
		t.Error("did not match")
	}

	dmg = client.Damage{Msg: "-l33t- MadJoker (Type 81 (C)) shot down -Tygri- xbladecz (F/A-18A)"}
	matched = d.parseDamage(dmg, 0)

	if matched == false {
		t.Error("did not match")
	}

}
