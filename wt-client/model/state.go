package model

type State struct {
	Volume   int     `json:"volume"`
	Music    []Music `json:"music"`
	Colldown uint16  `json:"cooldown"`
}

type StateRule struct {
	States         []string        `json:"states"`
	ConditionsBool map[string]bool `json:"conditions_bool"`
}

func (s *State) merge(st State) State {
	var dest State

	if st.Volume != 0 {
		dest.Volume = st.Volume
	} else {
		dest.Volume = s.Volume
	}
	if len(st.Music) > 0 {
		dest.Music = st.Music
	} else {
		dest.Music = s.Music
	}
	if st.Colldown > 0 {
		dest.Colldown = st.Colldown
	} else {
		dest.Colldown = s.Colldown
	}

	return dest
}
