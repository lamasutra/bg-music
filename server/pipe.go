package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/lamasutra/bg-music/config"
	"github.com/lamasutra/bg-music/player"
)

type PipeServer struct {
	state          *ServerState
	controlChannel chan string
	stateChannel   chan string
	eventChannel   chan string
}

func NewPipeServer() *PipeServer {
	return &PipeServer{}
}

func (p *PipeServer) Serve(conf *config.Config) {
	sleepTime := time.Millisecond * 50
	p.controlChannel = make(chan string)
	p.stateChannel = make(chan string)
	p.eventChannel = make(chan string)
	serverState := ServerState{
		config: conf,
		player: player.CreatePlayer(conf.PlayerType, conf.Volume),
	}
	p.state = &serverState

	defer p.state.player.Close()

	go handlePipeFile(p.controlChannel, "control.pipe", sleepTime)
	go handlePipeFile(p.stateChannel, "state.pipe", sleepTime)
	go handlePipeFile(p.eventChannel, "event.pipe", sleepTime)

	changeState("idle", p.state)

	for {
		select {
		case control := <-p.controlChannel:
			p.handleControl(control)
		case state := <-p.stateChannel:
			changeState(state, p.state)
		case event := <-p.eventChannel:
			triggerEvent(event, p.state)
		default:
			time.Sleep(sleepTime)
		}
		// fmt.Println("tick")
	}
}

func (p *PipeServer) Close() {
	close(p.controlChannel)
	close(p.stateChannel)
	close(p.eventChannel)
}

func handlePipeFile(ch chan string, filename string, sleepTime time.Duration) {
	var buffer []byte

	_, err := os.Stat(filename)
	if err != nil {
		syscall.Mkfifo(filename, 0666)
	}

	pipeFile, err := os.Open(filename)
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

// controls

func (p *PipeServer) loadConfig(data *LoadConfigData) {
	fmt.Println("loading config:", data)
	p.state.config.Events = data.Events
	p.state.config.States = data.States

	str, _ := json.MarshalIndent(p.state.config, "", "  ")
	fmt.Println("loaded:", string(str))
}

func (p *PipeServer) handleControl(control string) error {
	req := &Request{}
	err := json.Unmarshal([]byte(control), req)
	if err != nil {
		fmt.Println("error:", err, req)
		return err
	}
	if req == nil {
		fmt.Println("invalid request")
		return errors.New("invalid request")
	}
	fmt.Println("Received control:", req.Action)

	switch req.Action {
	case "load":
		loadData, ok := req.Data.(LoadConfigData)
		if !ok {
			str, _ := json.MarshalIndent(req.Data, "", "  ")
			fmt.Println("invalid data in request", string(str))
			return errors.New("invalid data in request")
		}
		p.loadConfig(&loadData)
	default:
		return errors.New("unknown action")
	}

	return nil
}
