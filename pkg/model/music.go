package model

type Music struct {
	Volume    int    `json:"volume"`
	Path      string `json:"path"`
	Skip      int    `json:"skip"`
	EndBefore int    `json:"endBefore"`
}

type MusicMetadata struct {
	Artist   string  `json:"artist"`
	Title    string  `json:"title"`
	Album    string  `json:"album"`
	Year     int     `json:"year"`
	Duration float64 `json:"duration"`
}
