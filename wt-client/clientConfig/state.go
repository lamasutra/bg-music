package clientConfig

type State struct {
	Volume *int    `json:"volume"`
	Music  []Music `json:"music"`
}

type StateRule struct {
	States         []string        `json:"states"`
	ConditionsBool map[string]bool `json:"conditions_bool"`
}

// func (sr *StateRule) checkConditions(input *input.WtInput) bool {
// 	return false
// }
