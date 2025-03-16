package server

import (
	"errors"

	"github.com/lamasutra/bg-music/model"
	"github.com/lamasutra/bg-music/ui"
)

type ServerState struct {
	state  string
	config *model.Config
	player *model.Player
}

type Server interface {
	Serve(*model.Config, *model.Player)
	Close()
}

func CreateServer(serverType string) (Server, error) {
	switch serverType {
	case "pipe":
		return NewPipeServer(), nil
	case "http":
		return NewHttpServer(), nil
	}

	return nil, errors.New("unknown server type " + serverType)
}

func triggerEvent(event string, srv *ServerState) error {
	ui.Debug("Received event:", event)
	ev, err := srv.config.GetEvent(event)
	if err != nil {
		ui.Error(err)
		return err
	}
	sfx, err := ev.GetRandomSfx()
	if err != nil {
		ui.Error(err)
		return err
	}

	_, err = (*srv.player).PlaySfx(sfx, srv.config)
	if err != nil {
		ui.Error(err)
		return err
	}

	return nil
}

func changeState(state string, srv *ServerState) error {
	ui.Debug("Received state:", state)
	if srv.state == state {
		ui.Debug("already ", state)
		return nil
	}
	srv.state = state
	st, err := srv.config.GetState(state)
	if err != nil {
		ui.Error(err)
		return err
	}
	music, err := srv.config.GetRandomStateMusic(state)
	if err != nil {
		ui.Error(err)
		return err
	}
	if st.Volume != nil {
		// (*srv.player).SetVolume(uint8(*st.Volume))
		_, err = (*srv.player).PlayMusicAtVolume(music, srv.config, uint8(*st.Volume))
		if err != nil {
			ui.Error(err)
			return err
		}
	} else {
		_, err = (*srv.player).PlayMusic(music, srv.config)
		if err != nil {
			ui.Error(err)
			return err
		}
	}

	return nil
}

func changeMusic(state string, srv *ServerState) error {
	ui.Debug("changing music")
	music, err := srv.config.GetRandomStateMusic(state)
	if err != nil {
		ui.Error(err)
		return err
	}
	_, err = (*srv.player).PlayMusic(music, srv.config)

	return err
}
