package player

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lamasutra/bg-music/wt-client/clientConfig"
)

type PipePlayer struct {
	controlPipe *os.File
	statePipe   *os.File
	eventPipe   *os.File
}

type PipeRequest struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func (p *PipePlayer) SendEventStates(ec *clientConfig.EventStates) error {
	req := PipeRequest{}
	req.Action = "load"
	req.Data = ec

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = p.controlPipe.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}
	return nil
}

func (p *PipePlayer) SendState(state string) error {
	_, err := p.statePipe.WriteString(state + "\n")

	return err
}

func (p *PipePlayer) TriggerEvent(event string) error {
	_, err := p.eventPipe.WriteString(event + "\n")

	return err
}

func (p *PipePlayer) Init(c *clientConfig.Config) {
	var err error
	p.controlPipe, err = os.OpenFile("../control.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	p.statePipe, err = os.OpenFile("../state.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	p.eventPipe, err = os.OpenFile("../event.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defaultTheme := c.Themes["default"]

	ec := clientConfig.EventStates{
		Events: defaultTheme.Events,
		States: defaultTheme.States,
	}

	fmt.Println("Sending default theme events and states...")
	err = p.SendEventStates(&ec)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *PipePlayer) Close() {
	if p.controlPipe != nil {
		p.controlPipe.Close()
	}
	if p.statePipe != nil {
		p.statePipe.Close()
	}
	if p.eventPipe != nil {
		p.eventPipe.Close()
	}
}
