package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bigbangit.com/event-music/config"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
)

func getSongPath(song *config.Song, c *config.Config) string {
	return c.Path + "/" + song.Path
}

func getFileExtension(path string) (string, error) {
	base := filepath.Base(strings.ToLower(path))
	extIndex := strings.LastIndex(base, ".")
	if extIndex == -1 {
		return "", fmt.Errorf("path does not contain a dot: %s", path)
	}

	return base[extIndex:], nil
}

func openSong(song *config.Song, c *config.Config) (beep.StreamSeekCloser, beep.Format, error) {
	path := getSongPath(song, c)
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

func playSong(song *config.Song, c *config.Config) (beep.StreamSeekCloser, error) {
	streamer, format, err := openSong(song, c)

	if err != nil {
		return nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	fmt.Printf("playing %v\n", song.Path)

	speaker.Play(streamer)

	return streamer, nil
}
