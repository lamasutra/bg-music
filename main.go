package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	// "bigbangit.com/event-music/config"
	// "bigbangit.com/event-music/model"
	"github.com/bg-music/config"
	"github.com/bg-music/model"
	"github.com/bg-music/player"
	"github.com/gopxl/beep/v2"
)

type StateMachine struct {
	currentState model.State
}

type cmdArgs struct {
	player *string
	server *string
	run    *bool
}

func main() {
	cmdArgs := registerFlags()
	if cmdArgs == nil {
		return
	}
	conf, err := config.Read("config.json")
	if err != nil {
		panic("Cannot read config.json")
	}

	fmt.Println("Running as", *cmdArgs.player, *cmdArgs.server)
	player := player.CreatePlayer(*cmdArgs.player, conf.Volume)
	music := model.Music{Path: "crusader/1/08 Track 8.mp3"}
	_, err = player.PlayMusic(&music, conf)
	if err != nil {
		fmt.Println(err)
		panic("cannot player music")
	}
	defer player.Close()
	for {
		time.Sleep(time.Millisecond * 100)
	}
}

func mainBak() {
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

	// var currentSong *config.Song
	var streamer beep.StreamSeekCloser
	eventChannel = make(chan string)

	// go handlePipeFile(eventChannel, sleepTime)

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
			// currentSong, _ = conf.GetRandomSong(currentEvent)
			// streamer, err = playSong(currentSong, conf)
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
				// currentSong, _ := conf.GetRandomSong(currentEvent)
				// streamer, err = playSong(currentSong, conf)
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

func registerFlags() *cmdArgs {
	var args cmdArgs
	args.player = flag.String("player", "local", "The music player type, local is default")
	args.server = flag.String("server", "pipe", "The server type, pipe is default")
	args.run = flag.Bool("run", true, "Run the app")

	// Use a flag with usage function as its value
	helpFlag := flag.Bool("h", false, usage())
	versionFlag := flag.Bool("v", false, "")
	flag.Parse()

	if flag.NFlag() == 0 || *helpFlag {
		fmt.Println(usage())
		return nil
	} else if *versionFlag {
		fmt.Println("version: poc")
		return nil
	}

	return &args
}

func usage() string {
	return `
Usage:
  -h|--help   Show this message and exit
  -v          Print version information

Flags:
  player        The music player type (default: "local")
  server        The server type, pipe is default (default: "pipe)
  run			Run the app
`
}
