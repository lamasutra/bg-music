package player

import (
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/lamasutra/bg-music/ui"
)

func TestSequencerStream(t *testing.T) {
	t.Log("testing sequencer")

	return

	seq1 := NewBeepSequencer(3, "seq1")
	seq2 := NewBeepSequencer(3, "seq2")
	ui.CreateUI("cli")

	format := beep.Format{
		SampleRate: 44100,
	}
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	speaker.Play(&seq1, &seq2)

	for {
		time.Sleep(time.Second)
	}

	beep.Seq()
}
