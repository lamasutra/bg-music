package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bg-music/config"
	"github.com/bg-music/model"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
)

type LocalPlayer struct {
	currentMusic  *model.Music
	volumePercent uint8
	volumeEffect  *effects.Volume
	streamer      beep.StreamSeekCloser
}

func (p *LocalPlayer) getSfxPath(sfx *model.Sfx, c *config.Config) string {
	return c.Path + "/" + sfx.Path
}

func (p *LocalPlayer) getMusicPath(music *model.Music, c *config.Config) string {
	return c.Path + "/" + music.Path
}

func (p *LocalPlayer) getFileExtension(path string) (string, error) {
	base := filepath.Base(strings.ToLower(path))
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return "", fmt.Errorf("path does not contain a dot: %s", path)
	}

	return base[extIndex:], nil
}

func (p *LocalPlayer) openFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
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

func (p *LocalPlayer) play(path string) (beep.StreamSeekCloser, error) {
	streamer, format, err := p.openFile(path)

	if err != nil {
		return nil, err
	}

	p.streamer = streamer
	p.volumeEffect = &effects.Volume{
		Base:     2,
		Silent:   false,
		Streamer: streamer,
	}
	p.SetVolume(p.volumePercent)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(p.volumeEffect)

	return streamer, nil
}

func (p *LocalPlayer) SetVolume(volume uint8) {
	fmt.Println("volume:", volume)

	if volume == 0 {
		p.volumeEffect.Silent = true
	} else {
		p.volumeEffect.Silent = false
		realVolume := float64(volume)/32 - 2 // math.Round(float64(volume-150) / 64)
		p.volumeEffect.Volume = realVolume
	}
}

func (p *LocalPlayer) PlayMusic(music *model.Music, c *config.Config) (beep.StreamSeekCloser, error) {
	path := p.getMusicPath(music, c)

	p.currentMusic = music

	fmt.Printf("playing music %v\n", path)

	return p.play(path)
}

func (p *LocalPlayer) PlaySfx(sfx *model.Sfx, c *config.Config) (beep.StreamSeekCloser, error) {
	path := p.getSfxPath(sfx, c)

	fmt.Printf("playing sfx %v\n", path)

	return p.play(path)
}

func (p *LocalPlayer) Close() {
	p.streamer.Close()
}
