package server

import (
	"errors"
	"fmt"

	"github.com/lamasutra/bg-music/config"
	"github.com/lamasutra/bg-music/player"
)

type ServerState struct {
	state  string
	config *config.Config
	player player.Player
}

type Server interface {
	Serve(*config.Config)
	Close()
}

func CreateServer(serverType string) (Server, error) {
	switch serverType {
	case "pipe":
		return NewPipeServer(), nil
	}

	return nil, errors.New("unknown server type " + serverType)
}

func triggerEvent(event string, srv *ServerState) error {
	fmt.Println("Received event:", event)
	et, err := srv.config.GetEvent(event)
	if err != nil {
		fmt.Println(err)
		return err
	}
	sfx, err := srv.config.GetRandomEventSfx(event)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if et.Volume != nil {
		srv.player.SetVolume(uint8(*et.Volume))
	}

	_, err = srv.player.PlaySfx(sfx, srv.config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func changeState(state string, srv *ServerState) error {
	fmt.Println("Received state:", state)
	srv.state = state
	st, err := srv.config.GetState(state)
	if err != nil {
		fmt.Println(err)
		return err
	}
	music, err := srv.config.GetRandomStateMusic(state)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if st.Volume != nil {
		srv.player.SetVolume(uint8(*st.Volume))
	}

	_, err = srv.player.PlayMusic(music, srv.config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func changeMusic(state string, srv *ServerState) error {
	fmt.Println("changing music")
	music, err := srv.config.GetRandomStateMusic(state)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = srv.player.PlayMusic(music, srv.config)

	return err
}
