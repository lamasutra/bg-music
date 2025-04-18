package model

type Event struct {
	Volume   int      `json:"volume"`
	Sfx      []Sfx    `json:"sfx"`
	Sentence []string `json:"sentence"`
}

func (e *Event) merge(ev Event) Event {
	var dest Event

	if ev.Volume != 0 {
		dest.Volume = ev.Volume
	} else {
		dest.Volume = e.Volume
	}
	if len(ev.Sfx) > 0 {
		dest.Sfx = ev.Sfx
	} else {
		dest.Sfx = e.Sfx
	}
	if len(ev.Sentence) > 0 {
		dest.Sentence = ev.Sentence
	} else {
		dest.Sentence = e.Sentence
	}

	return dest
}
