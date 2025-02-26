package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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

func main() {
	var newEvent, currentEvent string
	var eventChannel chan string

	conf, err := config.Read("config.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = os.Stat("event.pipe")
	if os.IsNotExist(err) {
		err = syscall.Mkfifo("event.pipe", 0666)
		if err != nil {
			fmt.Println("Make named pipe file error:", err)
			return
		}
	} else if err != nil {
		fmt.Println("Error checking file:", err)
	}

	sleepTime := time.Millisecond * 50

	var currentSong *config.Song
	var streamer beep.StreamSeekCloser
	eventChannel = make(chan string)

	go handlePipeFile(eventChannel, sleepTime)

	// @todo, make configurable
	newEvent = "idle"

	for {
		// fmt.Println(streamer)
		// if streamer != nil {
		// 	pos := uint64(streamer.Position())
		// 	length := uint64(streamer.Len())
		// 	posPerc := math.Round(float64(pos) / float64(length) * 100)
		// 	fmt.Println("position", pos, "of", length, posPerc, "%")
		// }

		if streamer != nil && streamer.Position() >= streamer.Len() {
			currentSong, _ = conf.GetRandomSong(currentEvent)
			streamer, err = playSong(currentSong, conf)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		select {
		case newEvent = <-eventChannel:
			// Process result from the channel when it arrives
			fmt.Println("Event:", newEvent)
		case <-time.After(sleepTime):
			// Timeout if no data is received in the specified time
			// fmt.Println("Timeout, no data received")
		}
		if len(newEvent) > 0 {
			if newEvent != currentEvent && conf.EventExists(newEvent) {
				fmt.Println("new event", newEvent)
				if streamer != nil {
					streamer.Close()
				}
				currentEvent = newEvent
				currentSong, _ := conf.GetRandomSong(currentEvent)
				streamer, err = playSong(currentSong, conf)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer streamer.Close()
			}
			newEvent = ""
			time.Sleep(sleepTime)
		}
	}
}

func handlePipeFile(ch chan string, sleepTime time.Duration) {
	var buffer []byte
	pipeFile, err := os.OpenFile("event.pipe", os.O_CREATE|os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println(err)
		return
	}
	pipeFileReader := bufio.NewReader(pipeFile)

	for {
		buffer, err = pipeFileReader.ReadBytes('\n')
		if err != nil {
			// fmt.Println("pipe error", err)
		} else {
			ch <- strings.TrimRight(string(buffer), "\n")
		}
		time.Sleep(sleepTime)
	}
}
