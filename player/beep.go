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
	"github.com/lamasutra/bg-music/config"
	"github.com/lamasutra/bg-music/model"
)

type BeepPlayer struct {
	initialized   bool
	currentMusic  *model.Music
	volumePercent uint8
	volumeEffect  *effects.Volume
	musicStreamer beep.StreamSeekCloser
	sfxStreamer   beep.StreamSeekCloser
	mixer         *beep.Mixer
	sampleRate    beep.SampleRate
	musicEnded    *chan (bool)
	crossfadeNum  int
}

func (p *BeepPlayer) getSfxPath(sfx *model.Sfx, c *config.Config) string {
	return c.Path + "/" + sfx.Path
}

func (p *BeepPlayer) getMusicPath(music *model.Music, c *config.Config) string {
	return c.Path + "/" + music.Path
}

func (p *BeepPlayer) getFileExtension(path string) (string, error) {
	base := filepath.Base(strings.ToLower(path))
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return "", fmt.Errorf("path does not contain a dot: %s", path)
	}

	return base[extIndex:], nil
}

func (p *BeepPlayer) openFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
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

func (p *BeepPlayer) play(streamer beep.Streamer, format beep.Format) error {
	fmt.Println("format:", format)
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

func (p *BeepPlayer) SetVolume(volume uint8) {
	p.volumePercent = volume

	if p.volumeEffect != nil {
		setVolume(p.volumeEffect, volume)
	}
}

func (p *BeepPlayer) PlayMusic(music *model.Music, c *config.Config) (beep.StreamSeekCloser, error) {
	path := p.getMusicPath(music, c)

	p.currentMusic = music

	streamer, format, err := p.openFile(path)

	if err != nil {
		return nil, err
	}

	if p.musicStreamer != nil {
		/* fallback
		p.mixer.Clear()
		p.musicStreamer = streamer
		p.SetVolume(p.volumePercent)
		p.mixer.Add(p.volumeEffect)
		*/
		p.volumeEffect = p.crossfadeCurrent(&streamer)
		p.musicStreamer = streamer

		// crossfade
		// cs1 := beep.Take(p.sampleRate.N(time.Second), p.musicStreamer)
		// cse1 := effects.Transition(
		// 	cs1,
		// 	p.sampleRate.N(time.Second),
		// 	p.currentGain()*float64(p.volumePercent)/100,
		// 	0.0,
		// 	effects.TransitionEqualPower,
		// )
		// // beep.Seq()
		// p.mixer.Clear()
		// time.Sleep(time.Second)
		// p.mixer.Add(cse1)
		// p.musicStreamer = streamer
		// p.SetVolume(p.volumePercent)
		// p.mixer.Add(p.volumeEffect)
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

	go p.watchStreamEnds()

	fmt.Printf("playing music %v\n", path)

	return streamer, err
}

func (p *BeepPlayer) GetMusicEndedChan() *chan (bool) {
	return p.musicEnded
}

func (p *BeepPlayer) PlaySfx(sfx *model.Sfx, c *config.Config) (beep.StreamSeekCloser, error) {
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

	fmt.Printf("playing sfx %v\n", path)

	return streamer, err
}

func (p *BeepPlayer) Init() {
	p.mixer = &beep.Mixer{}
	p.mixer.KeepAlive(true)
}

func (p *BeepPlayer) Close() {
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
		close(*p.musicEnded)
	}
}

func (p *BeepPlayer) watchStreamEnds() {
	sleepTime := time.Millisecond * 50
	for {
		if (p.musicStreamer.Position() - p.crossfadeNum) >= p.musicStreamer.Len() {
			fmt.Println("\rmusic " + p.currentMusic.Path + " ending")
			*p.musicEnded <- true
			return
		}
		fmt.Printf("\rpos: %ds of %ds", p.musicStreamer.Position()/int(p.sampleRate), p.musicStreamer.Len()/int(p.sampleRate))
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
	}

	return volumeEffect
}

func (p *BeepPlayer) crossfadeCurrent(streamer *beep.StreamSeekCloser) *effects.Volume {
	fmt.Println("crossfading", p.crossfadeNum)
	currentSample := beep.Take(p.crossfadeNum, p.volumeEffect)
	newVolumeEffect := wrapStreamerByVolumeEffect(streamer)
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
