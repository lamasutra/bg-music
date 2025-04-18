package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

type PipeServer struct {
	state             *ServerState
	controlChannel    chan string
	stateChannel      chan string
	eventChannel      chan string
	musicEndedChannel chan bool
}

func NewPipeServer() *PipeServer {
	return &PipeServer{}
}

func (p *PipeServer) Serve(conf *model.Config, player audio.Player) {
	sleepTime := time.Millisecond * 100
	p.controlChannel = make(chan string)
	p.stateChannel = make(chan string)
	p.eventChannel = make(chan string)
	p.musicEndedChannel = make(chan bool)
	serverState := ServerState{
		config: conf,
		player: player,
	}
	// player.CreatePlayer(conf.PlayerType, conf.Volume, &p.musicEndedChannel),
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
		case <-serverState.player.GetMusicEndedChan():
			changeMusic(p.state.state, p.state)
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
		logger.Debug(err)
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

func (p *PipeServer) loadConfig(data *LoadData) {
	// fmt.Println("loading config:", data)
	p.state.config.Events = data.Events
	p.state.config.States = data.States

	str, _ := json.MarshalIndent(p.state.config, "", "  ")
	logger.Debug(string(str))

	//
	// fmt.Println("loaded:", string(str))
}

func (p *PipeServer) handleControl(control string) error {
	req := &Request{}
	err := json.Unmarshal([]byte(control), req)
	if err != nil {
		logger.Error(err, req)
		return err
	}
	logger.Debug("Received control:", req.Action)

	switch req.Action {
	case "load":
		loadRequest := &LoadRequest{}
		err := json.Unmarshal([]byte(control), loadRequest)
		if err != nil {
			logger.Error("data", err)
			return err
		}
		p.loadConfig(&loadRequest.Data)
	case "next":
		changeMusic(p.state.state, p.state)
	default:
		return errors.New("unknown action")
	}

	return nil
}
