package model

type Event struct {
	Volume int     `json:"volume"`
	Music  []Music `json:"music"`
	Sfx    []Sfx   `json:"sfx"`
}
