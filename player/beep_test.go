package player

import (
	"testing"
	"time"

	"github.com/lamasutra/bg-music/model"
	"github.com/lamasutra/bg-music/ui"
)

var conf = model.Config{
	PlayerType: "beep",
	ServerType: "http",
	Volume:     100,
	Path:       "../music",
}

var music = model.Music{
	Path: "crusader/2/01 Track 1.mp3",
}
var music2 = model.Music{
	Path: "crusader/2/02 Track 2.mp3",
}

var sfx = model.Sfx{
	Volume: 30,
	Path:   "speech/tfx/kill.flac",
}

var sam = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/kill.sam",
		SampleRate: 8000,
	},
	Meaning: "kill",
}

var hostiles = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/hostiles.raw",
		SampleRate: 8000,
	},
	Meaning: "hostiles",
}

var s1 = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/1.sam",
		SampleRate: 8000,
	},
	Meaning: "1",
}
var s2 = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/2.raw",
		SampleRate: 8000,
	},
	Meaning: "2",
}
var s5 = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/5.sam",
		SampleRate: 8000,
	},
	Meaning: "5",
}
var s50 = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/50.sam",
		SampleRate: 8000,
	},
	Meaning: "50",
}

var s100 = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/100.sam",
		SampleRate: 8000,
	},
	Meaning: "100",
}

var degrees = model.Speech{
	Sfx: model.Sfx{
		Volume:     50,
		Path:       "speech/tfx/degrees.sam",
		SampleRate: 8000,
	},
	Meaning: "degrees",
}

var awacs1 = model.Speech{
	Sfx: model.Sfx{
		Volume:     100,
		Path:       "speech/tfx/awacs_confirms_inbound_hostiles_you_have_permission_to_fire.sam",
		SampleRate: 8000,
	},
	Meaning: "awacs1",
}

func TestCrossfade(t *testing.T) {
	t.Log("testing crossfade")

	ui.CreateUI("cli")
	streamer1, _, _ := openFile(getMusicPath(&music, &conf))
	streamer2, _, _ := openFile(getMusicPath(&music2, &conf))

	streamer1.Seek(10000)
	streamer2.Seek(10000)
	// crossfaded := crossfade(streamer1, streamer2, 1)

	samples := make([][2]float64, 512)

	// crossfaded.Stream(samples)

	ui.Debug(samples)

	streamer1.Stream(samples)
	ui.Debug(samples)
}

func TestSpeech(t *testing.T) {
	t.Log("testing speech")
	return
	// event := model.Event{
	// 	Volume: 100,
	// }
	// event.Sfx = append(event.Sfx, sfx)

	// conf.Events["test"] = event

	ui.CreateUI("cli")
	b := CreatePlayer("beep")

	// (*b).SetVolume(50)

	// sentence := []model.Speech{hostiles, s2, s100, s50, s2, degrees, awacs1}

	_, err := (*b).PlayMusic(&music, &conf)
	if err != nil {
		t.Error(err)
	}

	// (*b).Speak([]model.Speech{awacs1}, &conf)

	// (*b).Speak(&sentence, &conf)
	// (*b).PlaySfx(&sfx, &conf)

	// speaker.Play(&hostiles)

	i := 0
	for {
		time.Sleep(time.Second)
		i++
		if i == 3 {
			(*b).PlayMusic(&music2, &conf)
		}
	}

}
