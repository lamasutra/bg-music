package model

type Vehicle struct {
	Title  string `json:"title"`
	Theme  string `json:"theme"`
	Type   string `json:"type"`
	Volume int    `json:"volume"`
}

func (v *Vehicle) merge(to Vehicle) *Vehicle {
	if v.Theme != "" {
		to.Theme = v.Theme
	}
	if v.Title != "" {
		to.Title = v.Title
	}
	if v.Volume != 0 {
		to.Volume = v.Volume
	}
	if v.Type != "" {
		to.Type = v.Type
	}

	return &to
}
