package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"bigbangit.com/event-music/config"
	"bigbangit.com/event-music/model"
	"github.com/gopxl/beep/v2"
)

type StateMachine struct {
	currentState model.State
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
