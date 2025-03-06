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

type beepState struct {
	initialized   bool
	currentMusic  *model.Music
	volumePercent uint8
	volumeEffect  *effects.Volume
	musicStreamer beep.StreamSeekCloser
	sfxStreamer   beep.StreamSeekCloser
	mixer         *beep.Mixer
	sampleRate    beep.SampleRate
	crossfadeNum  int
	musicEnded    chan (bool)
	stopWatchEnd  chan (bool)
}

func (p *beepState) getSfxPath(sfx *model.Sfx, c *model.Config) string {
	return c.Path + "/" + sfx.Path
}

func (p *beepState) getMusicPath(music *model.Music, c *model.Config) string {
	return c.Path + "/" + music.Path
}

func (p *beepState) getFileExtension(path string) (string, error) {
	base := filepath.Base(strings.ToLower(path))
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return "", fmt.Errorf("path does not contain a dot: %s", path)
	}

	return base[extIndex:], nil
}

func (p *beepState) openFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	ext, err := p.getFileExtension(path)

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

func (p *beepState) play(streamer beep.Streamer, format beep.Format) error {
	ui.Debug("format:", format)
	if !p.initialized {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			return err
		}
		p.sampleRate = format.SampleRate
		speaker.Play(p.mixer)
		p.initialized = true
	}
	// if format.SampleRate != p.sampleRate {
	// 	resampler := beep.Resample(1, format.SampleRate, p.sampleRate, streamer)
	// 	stream = resampler.
	// }

	p.mixer.Add(streamer)

	return nil
}

func (p *beepState) SetVolume(volume uint8) {
	p.volumePercent = volume

	if p.volumeEffect != nil {
		setVolume(p.volumeEffect, volume)
	}
}

func (p *beepState) PlayMusic(music *model.Music, c *model.Config) (beep.StreamSeekCloser, error) {
	path := p.getMusicPath(music, c)

	p.currentMusic = music

	streamer, format, err := p.openFile(path)

	if err != nil {
		return nil, err
	}

	if p.musicStreamer != nil {
		p.volumeEffect = p.crossfadeCurrent(&streamer)
		p.musicStreamer = streamer
	} else {
		p.volumeEffect = wrapStreamerByVolumeEffect(&streamer)
		p.musicStreamer = streamer
		p.SetVolume(p.volumePercent)

		err = p.play(p.volumeEffect, format)
		if err != nil {
			return streamer, err
		}

		p.crossfadeNum = p.sampleRate.N(time.Second / 2)
	}

	// fmt.Println("str len", streamer.Len(), p.crossfadeNum)
	// beep.Take()streamer.

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
	if p.musicStreamer.Len() == 0 {
		ui.Error("cm is empty")
		return 0.0
	}

	return float64(p.musicStreamer.Position()) / float64(p.musicStreamer.Len())
}

func (p *beepState) PlaySfx(sfx *model.Sfx, c *model.Config) (beep.StreamSeekCloser, error) {
	path := p.getSfxPath(sfx, c)

	streamer, format, err := p.openFile(path)
	if err != nil {
		return nil, err
	}

	p.sfxStreamer = streamer

	err = p.play(streamer, format)
	if err != nil {
		return nil, err
	}

	ui.Debug(fmt.Sprintf("playing sfx %v\n", path))

	return streamer, err
}

func (p *beepState) Init() {
	p.mixer = &beep.Mixer{}
	p.mixer.KeepAlive(true)
	p.musicEnded = make(chan bool)
	go p.watchMusicStream()
}

func (p *beepState) Close() {
	if p.musicStreamer != nil {
		p.musicStreamer.Close()
	}

	if p.sfxStreamer != nil {
		p.sfxStreamer.Close()
	}

	if p.mixer != nil {
		p.mixer.Clear()
	}

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
			if (p.musicStreamer.Position() + p.crossfadeNum) >= p.musicStreamer.Len() {
				ui.Debug("music", p.currentMusic.Path, "ending", "mem", &p.currentMusic)
				p.musicEnded <- true
			}
		}
		ui.SetCurrentMusicProgress(p.GetCurrentMusicProgress())
		ui.Debug("p", p.GetCurrentMusicProgress(), "\n")
		time.Sleep(sleepTime)
	}
}

func wrapStreamerByVolumeEffect(streamer *beep.StreamSeekCloser) *effects.Volume {
	volumeEffect := effects.Volume{
		Base:     2,
		Silent:   false,
		Streamer: *streamer,
	}

	return &volumeEffect
}

func setVolume(volumeEffect *effects.Volume, volumePercent uint8) *effects.Volume {
	if volumePercent == 0 {
		volumeEffect.Silent = true
	} else {
		volumeEffect.Silent = false
		realVolume := float64(volumePercent)/20 - 5
		volumeEffect.Volume = realVolume
		ui.Debug("setVolume on", volumeEffect, "to", realVolume, volumePercent)
	}

	ui.SetCurrentVolume(float64(volumePercent) / 100)

	return volumeEffect
}

func (p *beepState) crossfadeCurrent(streamer *beep.StreamSeekCloser) *effects.Volume {
	ui.Debug("crossfading", p.crossfadeNum)
	currentSample := beep.Take(p.crossfadeNum, p.volumeEffect)
	newVolumeEffect := wrapStreamerByVolumeEffect(streamer)
	setVolume(newVolumeEffect, p.volumePercent)
	newSample := beep.Take(p.crossfadeNum, newVolumeEffect)
	mixed := crossfade(&currentSample, &newSample, p.crossfadeNum)

	(*streamer).Seek(p.crossfadeNum)
	seq := beep.Seq(*mixed, *streamer)

	p.mixer.Clear()
	p.mixer.Add(seq)

	return newVolumeEffect
}

func crossfade(stream1 *beep.Streamer, stream2 *beep.Streamer, length int) *beep.Streamer {
	trans1 := effects.Transition(*stream1, length, 1.0, 0.0, effects.TransitionEqualPower)
	trans2 := effects.Transition(*stream2, length, 0.0, 1.0, effects.TransitionEqualPower)
	mixed := beep.Take(length, beep.Mix(trans1, trans2))

	return &mixed
}
