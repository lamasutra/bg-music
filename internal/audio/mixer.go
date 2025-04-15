package audio

import "github.com/lamasutra/bg-music/pkg/logger"

type mixer struct {
	sequencers *[]sequencer
}

func NewBeepMixer(sequencers *[]sequencer) *mixer {
	return &mixer{
		sequencers: sequencers,
	}
}

// Add adds Streamers to the Mixer.
// func (m *mixer) Add(s ...*sequencer) {
// 	logger.Debug("mixer appending", s)
// 	m.sequencers = append(m.sequencers, s...)
// }

func (m *mixer) Stream(samples [][2]float64) (n int, ok bool) {
	logger.Debug("mixer stream", len(samples))
	if len((*m.sequencers)) == 0 {
		// return 0, true
		return len(samples), true
	}

	var tmp [512][2]float64

	for len(samples) > 0 {
		toStream := min(len(tmp), len(samples))

		// Clear the samples
		clear(samples[:toStream])

		snMax := 0
		for si := 0; si < len((*m.sequencers)); si++ {
			// Mix the stream
			logger.Debug("seq ", si)
			sn, sok := (*m.sequencers)[si].Stream(tmp[:toStream])
			for i := range tmp[:sn] {
				samples[i][0] += tmp[i][0]
				samples[i][1] += tmp[i][1]
			}
			if sn > snMax {
				snMax = sn
			}

			if sn < toStream || !sok {
				return n + snMax, true
			}
		}

		samples = samples[toStream:]
		n += toStream
	}

	return n, true
}

func (m *mixer) Err() error {
	return nil
}
