package model

import (
	"encoding/json"
	"testing"

	"github.com/lamasutra/bg-music/wt-client/internal/ui"
)

func TestMerge(t *testing.T) {
	var config Config
	config.Read("../wt-config.json")

	ui.CreateUI("cli")

	defaultTheme := *config.getTheme("default")
	adlib := *config.getTheme("tfx-adlib")

	if defaultTheme.Extend != "" {
		t.Errorf("default theme does not 	extends")
	}

	js, _ := json.MarshalIndent(adlib, "", "  ")
	if adlib.Extend == "" {
		t.Errorf("adlib theme does not extends %s", string(js))
	}

	events := adlib.Events
	if len(events) == 0 {
		t.Errorf("adlib theme does not haves events %s", string(js))
	}

	states := adlib.States
	if len(states) == 0 {
		t.Errorf("adlib theme does not haves states %s", string(js))
	}
}

func TestGetVehicleTheme(t *testing.T) {
	var config Config
	config.Read("../wt-config.json")

	ui.CreateUI("cli")

	config.GetThemeForVehicle(&Vehicle{
		Title: "air",
		Theme: "default",
	})

	vehicle := Vehicle{
		Title: "ef_2000_block_10",
		Theme: "tfx-adlib",
	}

	theme := config.GetThemeForVehicle(&vehicle)
	if theme.Extend == "" {
		t.Errorf("theme does not extend")
	}

	// js, _ := json.MarshalIndent(theme, "", "  ")
	// t.Errorf("incomplete %s", string(js))
}

func TestReferences(t *testing.T) {
	var config Config
	config.Read("../wt-config.json")

	ui.CreateUI("cli")

	vehicle1 := Vehicle{
		Title: "ef_2000_block_10",
		Theme: "tfx-adlib",
	}

	vehicle2 := Vehicle{
		Title: "air",
		Theme: "default",
	}

	theme1 := config.GetThemeForVehicle(&vehicle1)
	if theme1.Extend == "" {
		t.Errorf("theme does not extend")
	}

	theme2 := config.GetThemeForVehicle(&vehicle2)

	if &theme1 == &theme2 {
		t.Errorf("theme1 references theme2")
	}

	if theme1.States["clear"].Music[0].Path != "tfx-adlib/13-interdict.mp3" {
		t.Errorf("invalid track")
	}
	if theme2.States["clear"].Music[0].Path == "tfx-adlib/13-interdict.mp3" {
		t.Errorf("invalid track")
	}
}
