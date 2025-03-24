package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/lamasutra/bg-music/model"
	"github.com/lamasutra/bg-music/ui"
)

// func NewVerboseStreamer(streamer beep.Streamer) verboseStream {
// 	return verboseStream{
// 		streamer,
// 	}
// }

// type verboseStream struct {
// 	streamer beep.Streamer
// }

// func (v verboseStream) Stream(samples [][2]float64) (n int, ok bool) {
// 	n, ok = v.streamer.Stream(samples)

// 	fmt.Println("DEBUG: streaming", len(samples), n, ok)
// 	fmt.Println(samples)

// 	return n, ok
// }

// func (v verboseStream) Err() error {
// 	return v.streamer.Err()
// }

type sequencers struct {
	music    *sequencer
	sfx      *sequencer
	narrator *sequencer
}

type beepState struct {
	initialized bool

	currentMusic  *model.Music
	volumePercent uint8
	volumed       *effects.Volume
	musicStreamer *beep.StreamSeekCloser
	sfxStreamer   *beep.StreamSeekCloser
	mixer         *beep.Streamer // wrapped by volumeEffect
	sequencers    *sequencers
	format        beep.Format
	sampleRate    beep.SampleRate
	crossfadeNum  int
	musicEnded    chan (bool)
	stopWatchEnd  chan (bool)
}

var speechCache = make(map[string]effects.Volume, 0)

func getSfxPath(sfx *model.Sfx, c *model.Config) string {
	return c.Path + "/" + sfx.Path
}

func getMusicPath(music *model.Music, c *model.Config) string {
	return c.Path + "/" + music.Path
}

func getFileExtension(path string) (string, error) {
	base := filepath.Base(strings.ToLower(path))
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return "", fmt.Errorf("path does not contain a dot: %s", path)
	}

	return base[extIndex:], nil
}

func openFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	ext, err := getFileExtension(path)

	if err != nil {
		return nil, beep.Format{}, err
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("cannot read file %v", path)
	}

	switch ext {
	case ".mp3":
		// fmt.Println("mp3")
		return mp3.Decode(file)
	case ".flac":
		// fmt.Println("flac")
		return flac.Decode(file)
	case ".ogg":
		// fmt.Println("ogg/vorbis")
		return vorbis.Decode(file)
		// case "mid":
		// return midi.Decode(file)
	}

	return nil, beep.Format{}, fmt.Errorf("cannot decode file type %v", ext)
}

func (p *beepState) playMusic(streamer *beep.Streamer) {
	fmt.Printf("appending music %p\n", streamer)
	p.sequencers.music.Append(streamer)
}

func (p *beepState) playSfx(streamer *beep.Streamer) {
	fmt.Println("appending sfx")
	p.sequencers.sfx.Append(streamer)
}

func (p *beepState) playSpeech(streamer *beep.Streamer) {
	fmt.Println("appending speech")
	p.sequencers.narrator.Append(streamer)
}

// func (p *beepState) play(streamer beep.Streamer, format beep.Format) error {
// 	ui.Debug("format:", format)
// 	if !p.initialized {
// 		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
// 		if err != nil {
// 			return err
// 		}
// 		p.sampleRate = format.SampleRate
// 		speaker.Play(p.mixer)
// 		p.initialized = true
// 	}
// 	// if format.SampleRate != p.sampleRate {
// 	// 	resampler := beep.Resample(1, format.SampleRate, p.sampleRate, streamer)
// 	// 	stream = resampler.
// 	// }

// 	p.mixer.Add(streamer)

// 	return nil
// }

func (p *beepState) SetVolume(volume uint8) {
	p.volumePercent = volume
	setVolume(p.volumed, volume)
}

func (p *beepState) PlayMusic(music *model.Music, c *model.Config) (beep.StreamSeekCloser, error) {
	return p.PlayMusicAtVolume(music, c, p.volumePercent)
}

func (p *beepState) PlayMusicAtVolume(music *model.Music, c *model.Config, volume uint8) (beep.StreamSeekCloser, error) {
	path := getMusicPath(music, c)

	p.currentMusic = music

	streamer, format, err := openFile(path)

	if err != nil {
		return nil, err
	}

	fmt.Println("DEBUG", "PlayMusicAtVolume", music, volume, p.musicStreamer)

	var volumed effects.Volume
	if p.musicStreamer == nil {
		stream, ok := streamer.(beep.Streamer)
		if !ok {
			panic("invalid streamer type")
		}
		volumed = wrapSimpleStreamerByVolumeEffect(&stream)
		setVolume(&volumed, volume)
	} else {
		volumed = p.crossfadeCurrent(&streamer, volume)
	}

	p.sequencers.music.Append(volumed)
	p.musicStreamer = &streamer

	ui.Debug(fmt.Sprintf("playing music %v, duration: %vs", path, streamer.Len()/int(format.SampleRate)))
	ui.SetCurrentMusicTitle(path)

	return streamer, err
}

func (p *beepState) GetMusicEndedChan() chan (bool) {
	return p.musicEnded
}

func (p *beepState) GetCurrentMusic() *model.Music {
	return p.currentMusic
}

func (p *beepState) GetCurrentMusicProgress() float64 {
	if p.musicStreamer == nil {
		ui.Error("cm is nil")
		return 0.0
	}
	if (*p.musicStreamer).Len() == 0 {
		ui.Error("cm is empty")
		return 0.0
	}

	return float64((*p.musicStreamer).Position()) / float64((*p.musicStreamer).Len())
}

func (p *beepState) PlaySfx(sfx *model.Sfx, c *model.Config) (beep.StreamSeekCloser, error) {
	path := getSfxPath(sfx, c)

	streamer, format, err := openFile(path)
	if err != nil {
		return nil, err
	}

	p.sfxStreamer = &streamer

	ui.Debug(sfx)

	volumeSfxStreamer := wrapStreamerByVolumeEffect(&streamer)
	setVolume(&volumeSfxStreamer, sfx.Volume)

	var preparedStreamer beep.Streamer

	if format.SampleRate != p.format.SampleRate {
		preparedStreamer = beep.Resample(3, format.SampleRate, p.format.SampleRate, &volumeSfxStreamer)
	} else {
		preparedStreamer = &volumeSfxStreamer
	}

	p.playSfx(&preparedStreamer)

	ui.Debug("playing sfx ", path, sfx.Volume)

	return streamer, err
}

func (p *beepState) Speak(sentence *[]model.Speech, c *model.Config) {
	var streamer effects.Volume
	if len(*sentence) == 0 {
		return
	}
	for _, speech := range *sentence {
		err := speech.Prepare(c)
		if err != nil {
			ui.Error(speech.Meaning, "cannot prepare speech", err)
			continue
		}
		resampler := speech.Resample(p.sampleRate)
		streamer = wrapResamplerByVolumeEffect(resampler)
		setVolume(&streamer, speech.Volume)
		speechCache[speech.Meaning] = streamer
		p.sequencers.narrator.Append(&streamer)
	}
}

func (p *beepState) Play(s beep.Streamer) {
	speaker.Play(s)
}

func (p *beepState) Init() {
	p.musicEnded = make(chan bool)
	// p.mixer = NewBeepMixer()
	// p.mixer = beep.Mix(
	// 	(*p.sequencers).music,
	// 	(*p.sequencers).sfx,
	// 	(*p.sequencers).narrator,
	// )
	p.crossfadeNum = p.format.SampleRate.N(time.Second) / 2
	wrapped := wrapSimpleStreamerByVolumeEffect(p.mixer)
	p.volumed = &wrapped
	setVolume(p.volumed, 100)

	err := speaker.Init(p.format.SampleRate, p.format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	ui.Debug(fmt.Sprintf("beep player initialized, seqs: music=%p sfx=%p narrator=%p", &p.sequencers.music, &p.sequencers.sfx, &p.sequencers.narrator))
	p.sampleRate = p.format.SampleRate
	speaker.Play(p.volumed)
	// speaker.Play(p.mixer)
	p.initialized = true
	go p.watchMusicStream()
}

func (p *beepState) Close() {
	(*p.musicStreamer).Close()
	// p.sfxStreamer.Close()
	// p.mixer.Clear()

	if p.musicEnded != nil {
		close(p.musicEnded)
	}
	if p.stopWatchEnd != nil {
		p.stopWatchEnd <- true
		time.Sleep(time.Millisecond * 5)
		close(p.stopWatchEnd)
	}
}

func (p *beepState) watchMusicStream() {
	sleepTime := time.Millisecond * 100
	ui.Debug("entering watchStreamEnds")
	for {
		if p.musicStreamer == nil {
			ui.Debug("no stream yet")
			time.Sleep(sleepTime)
			continue
		}

		select {
		case <-p.stopWatchEnd:
			ui.Debug("exiting watchStreamEnds")
			return
		default:
			// ui.
			if ((*p.musicStreamer).Position() + p.crossfadeNum) >= (*p.musicStreamer).Len() {
				ui.Debug("music", p.currentMusic.Path, "ending", "mem", &p.currentMusic)
				p.musicEnded <- true
			}
		}
		ui.SetCurrentMusicProgress(p.GetCurrentMusicProgress())
		// ui.Debug("p", p.GetCurrentMusicProgress())
		time.Sleep(sleepTime)
	}
}

func wrapStreamerByVolumeEffect(streamer *beep.StreamSeekCloser) effects.Volume {
	return effects.Volume{
		Base:     2,
		Silent:   false,
		Streamer: *streamer,
	}
}

func wrapSimpleStreamerByVolumeEffect(streamer *beep.Streamer) effects.Volume {
	return effects.Volume{
		Base:     2,
		Silent:   false,
		Streamer: *streamer,
	}
}

func wrapResamplerByVolumeEffect(streamer *beep.Resampler) effects.Volume {
	return effects.Volume{
		Base:     2,
		Silent:   false,
		Streamer: streamer,
	}
}

func setVolume(volumeEffect *effects.Volume, volumePercent uint8) {
	if volumePercent == 0 {
		volumeEffect.Silent = true
	} else {
		volumeEffect.Silent = false
		realVolume := float64(volumePercent)/20 - 5
		volumeEffect.Volume = realVolume
		ui.Debug("setVolume on", volumeEffect, "to", realVolume, volumePercent)
	}

	ui.SetCurrentVolume(float64(volumePercent) / 100)
}

func (p *beepState) crossfadeCurrent(streamer *beep.StreamSeekCloser, newVolume uint8) effects.Volume {
	ui.Debug("crossfading ", p.crossfadeNum)
	ui.Debug(fmt.Sprintf("current: %p, new: %p", p.musicStreamer, streamer))
	currentSample := beep.Take(p.crossfadeNum, p.volumed)
	newVolumeEffect := wrapStreamerByVolumeEffect(streamer)
	setVolume(&newVolumeEffect, newVolume)
	newSample := beep.Take(p.crossfadeNum, &newVolumeEffect)
	mixed := crossfade(&currentSample, &newSample, p.crossfadeNum)

	(*p.sequencers.music).Append(&mixed)
	(*p.sequencers.music).Append(newVolumeEffect)

	(*p.musicStreamer).Seek((*p.musicStreamer).Len())
	(*streamer).Seek(p.crossfadeNum)

	return newVolumeEffect
}

// func (b *beepState) prepareSpeech(speech model.Speech) *beep.Streamer {
// 	var streamer *beep.Streamer
// 	// speechSampelRate := beep.SampleRate(speech.SampleRate)
// 	// if b.sampleRate != speechSampelRate {
// 	// 	beep.Resample(3, speechSampelRate, b.sampleRate)
// 	// }

// 	return streamer
// }

func crossfade(stream1 *beep.Streamer, stream2 *beep.Streamer, length int) beep.Streamer {
	trans1 := effects.Transition(*stream1, length, 1.0, 0.0, effects.TransitionEqualPower)
	trans2 := effects.Transition(*stream2, length, 0.0, 1.0, effects.TransitionEqualPower)
	mixed := beep.Take(length, beep.Mix(trans1, trans2))

	return mixed
}
