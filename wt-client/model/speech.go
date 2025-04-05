package model

type Speech struct {
	Sfx
	Meaning string `json:"meaning"`
	Format  int    `json:"format"`
	Skip    []int  `json:"skip"`
}

func (s *Speech) merge(st Speech) Speech {
	var dest Speech

	if st.Volume != 0 {
		dest.Volume = st.Volume
	} else {
		dest.Volume = s.Volume
	}

	if st.Format != 0 {
		dest.Format = st.Format
	} else {
		dest.Format = s.Format
	}

	if st.NumChannels != 0 {
		dest.NumChannels = st.NumChannels
	} else {
		dest.NumChannels = s.NumChannels
	}

	if st.SampleRate != 0 {
		dest.SampleRate = st.SampleRate
	} else {
		dest.SampleRate = s.SampleRate
	}

	if len(st.Skip) > 0 {
		dest.Skip = st.Skip
	} else {
		dest.Skip = s.Skip
	}

	if st.Meaning != "" {
		dest.Meaning = st.Meaning
	} else {
		dest.Meaning = s.Meaning
	}

	if st.Path != "" {
		dest.Path = st.Path
	} else {
		dest.Path = s.Path
	}

	return dest
}
