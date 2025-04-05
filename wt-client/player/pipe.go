package player

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/ui"
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

func (p *PipePlayer) SendEventStates(ec *model.BgPlayerConfig) error {
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

	ui.Debug("state sent", state)

	return err
}

func (p *PipePlayer) TriggerEvent(event string) error {
	_, err := p.eventPipe.WriteString(event + "\n")

	return err
}

func (p *PipePlayer) ChangeMusic() error {
	_, err := p.controlPipe.WriteString("next\n")

	return err
}

func (p *PipePlayer) Speak(string) error {
	ui.Debug("speak is not supported yet")

	return nil
}

func (p *PipePlayer) Init(c *model.Config) {
	ui.Debug("waiting for connection to bg player")
	var err error
	p.controlPipe, err = os.OpenFile("../control.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	ui.Debug("control pipe opened")
	p.statePipe, err = os.OpenFile("../state.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	ui.Debug("state pipe opened")
	p.eventPipe, err = os.OpenFile("../event.pipe", os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	ui.Debug("event pipe opened")

	defaultTheme := c.Themes["default"]

	ec := model.BgPlayerConfig{
		Events: defaultTheme.Events,
		States: defaultTheme.States,
	}

	ui.Debug("Sending default theme events and states...")
	err = p.SendEventStates(&ec)
	if err != nil {
		ui.Error(err)
	}
	time.Sleep(time.Millisecond * 50)
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
