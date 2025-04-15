package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lamasutra/bg-music/wt-client/internal/model"
	"github.com/lamasutra/bg-music/wt-client/internal/ui"
)

type HttPlayer struct {
	client *http.Client
	host   string
}

func (h *HttPlayer) SendEventStates(ec *model.BgPlayerConfig) error {
	data, err := json.Marshal(*ec)
	if err != nil {
		return err
	}

	resp, err := h.client.Post(h.host+"control/load", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: %s", string(body))
	}
	defer resp.Body.Close()

	return nil
}

// func (p *PipePlayer) ChangeVehicle(v *model.Vehicle) {

// }

func (h *HttPlayer) SendState(state string) error {
	req, err := http.NewRequest("PUT", h.host+"state/"+state, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	_, err = h.client.Do(req)
	if err != nil {
		return err
	}
	req.Body.Close()

	ui.Debug("state sent", state)

	return err
}

func (h *HttPlayer) TriggerEvent(event string) error {
	req, err := http.NewRequest("PUT", h.host+"event/"+event, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	_, err = h.client.Do(req)
	if err != nil {
		return err
	}
	req.Body.Close()

	ui.Debug("event triggered", event)

	return err
}

func (h *HttPlayer) ChangeMusic() error {
	req, err := http.NewRequest("POST", h.host+"control/next", bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	_, err = h.client.Do(req)
	if err != nil {
		return err
	}
	req.Body.Close()

	ui.Debug("changed music")

	return err
}

func (h *HttPlayer) Speak(sentence string) error {
	req, err := http.NewRequest("PUT", h.host+"speak?sentence="+sentence, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	_, err = h.client.Do(req)
	if err != nil {
		return err
	}
	req.Body.Close()

	ui.Debug("speak sent", sentence)

	return err
}

func (h *HttPlayer) Init(c *model.Config) {
	h.client = &http.Client{}
	h.host = c.BgPlayerHost
	ui.Debug("waiting for connection to bg player")

	for {
		resp, err := h.client.Get(c.BgPlayerHost)
		if err != nil {
			ui.Error("connection not ready")
		} else {
			resp.Body.Close()
			break
		}
		time.Sleep(time.Second)
	}

	defaultTheme := c.Themes["default"]

	ec := model.BgPlayerConfig{
		Events:  defaultTheme.Events,
		States:  defaultTheme.States,
		Narrate: defaultTheme.Narrate,
	}

	ui.Debug("Sending default theme events and states...")
	err := h.SendEventStates(&ec)
	if err != nil {
		ui.Error(err)
	}
	time.Sleep(time.Millisecond * 50)
}

func (h *HttPlayer) Close() {

}
