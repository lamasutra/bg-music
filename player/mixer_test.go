package player

import (
	"os"
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/lamasutra/bg-music/ui"
)

// var conf = model.Config{
// 	PlayerType: "beep",
// 	ServerType: "http",
// 	Volume:     100,
// 	Path:       "../music",
// }

func TestMixer(t *testing.T) {
	// return
	ui.CreateUI("cli")
	// CreatePlayer("beep")

	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	m := NewBeepMixer()
	s := NewBeepSequencer(10)
	x := NewBeepSequencer(10)
	file, _ := os.Open("../music/crusader/1/01 Track1.mp3")
	streamer, format, _ := mp3.Decode(file)
	s.Append(&streamer)

	// hostiles := model.Speech{
	// 	Sfx: model.Sfx{
	// 		Volume:     50,
	// 		Path:       "speech/tfx/hostiles.raw",
	// 		SampleRate: 8000,
	// 	},
	// 	Meaning: "hostiles",
	// }
	// vs := verboseStream{
	// streamer: &sp,
	// }

	// s.Append(&hostiles)

	m.Add(s)
	m.Add(x)

	// samples := make([][2]float64, 512)
	// n, ok := m.Stream(samples)
	// if !ok && n == 0 {
	// 	t.Error("mixer should stream")
	// }
	// (*b).Play(&sp)
	// sentence := []model.Speech{hostiles}

	// (*b).Speak(&sentence, &conf)
	speaker.Play(&m)
	for {
		time.Sleep(time.Second)
	}
}
