package audio

import (
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

func TestSequencerStream(t *testing.T) {
	t.Log("testing sequencer")

	return

	seq1 := NewBeepSequencer(3, "seq1")
	seq2 := NewBeepSequencer(3, "seq2")

	format := beep.Format{
		SampleRate: 44100,
	}
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	speaker.Play(&seq1, &seq2)

	i := 0
	for {
		time.Sleep(time.Second)
		if i == 5 {
			seq1.Append(&s1)
		}
		if i == 10 {
			seq1.Append(&s1)
		}

		i++
	}

}
