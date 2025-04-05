package model

type Music struct {
	Volume    int    `json:"volume"`
	Path      string `json:"path"`
	Skip      int    `json:"skip"`
	EndBefore int    `json:"endBefore"`
}
