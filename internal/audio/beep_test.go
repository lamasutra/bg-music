package audio

import (
	"fmt"
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

var conf = model.Config{
	PlayerType: "beep",
	ServerType: "http",
	Volume:     100,
	Path:       "../music",
}

var music1 = model.Music{
	Volume: 100,
	Path:   "crusader/2/01 Track 1.mp3",
}
var music2 = model.Music{
	Volume: 100,
	Path:   "crusader/2/02 Track 2.mp3",
}
var music3 = model.Music{
	Volume: 100,
	Path:   "crusader/2/03 Track 3.mp3",
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

func TestMusic(t *testing.T) {
	t.Log("testing music")

	// return

	b := CreatePlayer("beep")
	b.PlayMusic(&music1, &conf, true)
	i := 0
	for {
		time.Sleep(time.Second)
		i++
		if i == 3 {
			break
		}
	}

	b.PlayMusic(&music2, &conf, true)
	for {
		time.Sleep(time.Second)
		i++
		if i == 6 {
			break
		}
	}

	b.PlayMusic(&music1, &conf, true)
	for {
		time.Sleep(time.Second)
		i++
		if i == 9 {
			break
		}
	}

	logger.Debug("close")
	b.Close()
}

func TestCrossfade(t *testing.T) {
	t.Log("testing crossfade")

	return

	format := beep.Format{SampleRate: 44100}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	crossfadeNum := format.SampleRate.N(time.Second)

	streamer1, _, _ := openFile(conf.GetMusicPath(&music3))
	streamer2, _, _ := openFile(conf.GetMusicPath(&music2))

	s := NewBeepSequencer(16, "song")
	ss := NewBeepSequencer(32, "sfx")
	s.Append(streamer1)
	// ss.Append(streamer2)

	m := beep.Mix(&s, &ss)

	// streamer1.Seek(1024)
	// streamer2.Seek(0)
	// crossfaded := crossfade(streamer1, streamer2, crossfadeNum)
	speaker.Play(m)

	// samples := make([][2]float64, 512)

	// crossfaded.Stream(samples)

	// streamer1.Stream(samples)
	// ui.Debug(samples)

	// streamer1.Stream(samples)
	// ui.Debug(samples)
	i := 0
	for {
		fmt.Println(i)
		time.Sleep(time.Second)
		if i == 5 {
			crossfaded := crossfade(streamer1, streamer2, crossfadeNum)
			streamer1.Seek(streamer1.Len())
			s.Append(crossfaded)
			streamer2.Seek(crossfadeNum)
			s.Append(streamer2)
		}
		i++
	}
}

func TestSpeech(t *testing.T) {
	t.Log("testing speech")

	return

	// event := model.Event{
	// 	Volume: 100,
	// }
	// event.Sfx = append(event.Sfx, sfx)

	// conf.Events["test"] = event

	b := CreatePlayer("beep")
	// b.SetVolume(100)

	// sentence := []model.Speech{hostiles, s2, s100, s50, s2, degrees, awacs1}

	b.PlayMusic(&music1, &conf, true)

	// b.Speak([]model.Speech{awacs1}, &conf)

	// b.Speak(&sentence, &conf)
	// (*b).PlaySfx(&sfx, &conf)

	// speaker.Play(&hostiles)

	i := 0
	for {
		time.Sleep(time.Second)
		i++
		// if i == 3 {
		// 	b.PlayMusic(&music2, &conf)
		// }
	}

}
