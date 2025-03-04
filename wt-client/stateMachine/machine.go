package stateMachine

import (
	"errors"

	"github.com/lamasutra/bg-music/wt-client/clientConfig"
	"github.com/lamasutra/bg-music/wt-client/types"
)

type StateMachine struct {
	state string
	rules *map[string]clientConfig.StateRule
}

func New(state string, rules *map[string]clientConfig.StateRule) *StateMachine {
	return &StateMachine{
		state: state,
		rules: rules,
	}
}

func (sm *StateMachine) GetCurrentState() string {
	return sm.state
}

func (sm *StateMachine) GetNextState(input *types.WtInputMapBool) (string, error) {
	possibleStateCodes, err := sm.getPossibleStateCodes()
	// fmt.Println("possible states", possibleStateCodes)
	if err != nil {
		return "", err
	}
	for _, code := range *possibleStateCodes {
		rule, err := sm.getStateRule(code)
		if err != nil {
			return "", err
		}
		// fmt.Println("checking for state", code)

		if checkConditions(input, rule) {
			return code, nil
		}
	}

	return "", nil
}

func (sm *StateMachine) SetState(state string) {
	sm.state = state
}

func (sm *StateMachine) getPossibleStateCodes() (*[]string, error) {
	state, err := sm.getStateRule(sm.state)
	if err != nil {
		return nil, err
	}

	return &state.States, nil
}

func (sm *StateMachine) getStateRule(state string) (*clientConfig.StateRule, error) {
	rule, ok := (*sm.rules)[state]
	if !ok {
		return nil, errors.New("unknown state " + state)
	}

	return &rule, nil
}

func checkConditions(input *types.WtInputMapBool, rule *clientConfig.StateRule) bool {
	matches := true
	for key, val := range rule.ConditionsBool {
		// fmt.Println("  checking", key, "for", val)
		if (*input)[key] != val {
			// fmt.Println("    failed at", key, /"value shoud be", val, "is", (*input)[key])
			matches = false
			break
		}
	}

	return matches
}
