package player

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lamasutra/bg-music/wt-client/model"
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

func (p *PipePlayer) SendEventStates(ec *model.EventStates) error {
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

// func (p *PipePlayer) ChangeVehicle(v *model.Vehicle) {

// }

func (p *PipePlayer) SendState(state string) error {
	_, err := p.statePipe.WriteString(state + "\n")

	fmt.Println("state sent", state)

	return err
}

func (p *PipePlayer) TriggerEvent(event string) error {
	_, err := p.eventPipe.WriteString(event + "\n")

	return err
}

func (p *PipePlayer) Init(c *model.Config) {
	fmt.Println("waiting for connection to bg player")
	var err error
	p.controlPipe, err = os.OpenFile("../control.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("control pipe opened")
	p.statePipe, err = os.OpenFile("../state.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("state pipe opened")
	p.eventPipe, err = os.OpenFile("../event.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("event pipe opened")

	defaultTheme := c.Themes["default"]

	ec := model.EventStates{
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
