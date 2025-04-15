package model

type State struct {
	Volume         int     `json:"volume"`
	Music          []Music `json:"music"`
	Cooldown       int64   `json:"cooldown"`
	BreaksCooldown int8    `json:"breaks_cooldown"`
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
	if st.Cooldown > 0 {
		dest.Cooldown = st.Cooldown
	} else {
		dest.Cooldown = s.Cooldown
	}
	if st.BreaksCooldown != 0 {
		dest.BreaksCooldown = st.BreaksCooldown
	} else {
		dest.BreaksCooldown = s.BreaksCooldown
	}

	return dest
}
