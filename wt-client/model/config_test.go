package model

import (
	"encoding/json"
	"testing"
)

func TestMerge(t *testing.T) {
	var config Config
	config.Read("../wt-config.json")

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

	vehicle := Vehicle{
		Title: "EF2000",
		Theme: "tfx-adlib",
	}

	theme := config.GetThemeForVehicle(&vehicle)
	if theme.Extend == "" {
		t.Errorf("theme does not extend")
	}

	// js, _ := json.MarshalIndent(theme, "", "  ")
	// t.Errorf("incomplete %s", string(js))
}
