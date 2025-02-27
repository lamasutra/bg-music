package model

type Event struct {
	Volume *int  `json:"volume"`
	Sfx    []Sfx `json:"sfx"`
}
