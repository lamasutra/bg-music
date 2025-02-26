package model

type Music struct {
	Path      string `json:"path"`
	Skip      int    `json:"skip"`
	EndBefore int    `json:"endBefore"`
}
