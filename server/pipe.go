package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type PipeChannels struct {
	controlChannel chan string
	stateChannel   chan string
	eventChannel   chan string
}

func (p *PipeChannels) Serve() {
	sleepTime := time.Millisecond * 50
	p.controlChannel = make(chan string)
	p.stateChannel = make(chan string)
	p.eventChannel = make(chan string)

	go handlePipeFile(p.controlChannel, "control.pipe", sleepTime)
	go handlePipeFile(p.stateChannel, "state.pipe", sleepTime)
	go handlePipeFile(p.eventChannel, "event.pipe", sleepTime)

	for {
		select {
		case control := <-p.controlChannel:
			fmt.Println("Received control:", control)
		case state := <-p.stateChannel:
			fmt.Println("Received state:", state)
		case event := <-p.eventChannel:
			fmt.Println("Received event:", event)
		default:
			time.Sleep(sleepTime)
		}
	}

}

func handlePipeFile(ch chan string, filename string, sleepTime time.Duration) {
	var buffer []byte
	pipeFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println(err)
		return
	}
	pipeFileReader := bufio.NewReader(pipeFile)

	defer pipeFile.Close()

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
