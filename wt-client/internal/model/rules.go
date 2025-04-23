package model

import (
	"encoding/json"
	"os"

	"github.com/lamasutra/bg-music/pkg/logger"
)

type StateRules map[string]StateRule

type StateRule struct {
	States         []string        `json:"states"`
	ConditionsBool map[string]bool `json:"conditions_bool"`
	// ConditionsDistance map[string]uint64 `json:"conditions_bool"`
}

func (sr *StateRules) Read(path string) error {
	data, err := os.ReadFile(path)

	if err != nil {
		logger.Error("Cannot open rules file", path)
		return err
	}

	err = json.Unmarshal(data, &sr)
	if err != nil {
		logger.Error("Cannot decode json", err)
		return err
	}

	return nil
}
