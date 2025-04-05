package model

type Sfx struct {
	Volume      uint8  `json:"volume"`
	Path        string `json:"path"`
	SampleRate  int64  `json:"sample_rate"`
	NumChannels int    `json:"num_channels"`
}
