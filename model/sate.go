package model

type State struct {
	Volume *int     `json:"volume"`
	Music  []Music  `json:"music"`
	States []string `json:"states"`
}
