package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

type ServerState struct {
	state  string
	config *model.Config
	player audio.Player
}

type Server interface {
	Serve(*model.Config, audio.Player)
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
	logger.Debug("Received event:", event)
	ev, err := srv.config.GetEvent(event)
	if err != nil {
		logger.Error(err)
		return err
	}
	if len(ev.Sentence) > 0 {
		sentence, err := ev.GetRandomSentence()
		if err != nil {
			logger.Error(err)
			return err
		}

		err = speak(sentence, srv)
		if err != nil {
			logger.Error(err)
			return err
		}

	} else {
		sfx, err := ev.GetRandomSfx()
		if err != nil {
			logger.Error(err)
			return err
		}

		srv.player.PlaySfx(sfx, srv.config)
	}

	return nil
}

func speak(sentence string, srv *ServerState) error {
	seq := strings.Split(sentence, ",")
	logger.Debug(fmt.Sprintf("I will speak `%s` for you", sentence))
	narSeq := make([]model.Speech, len(seq))
	for i, key := range seq {
		val, ok := srv.config.Narrate[key]
		if !ok {
			return fmt.Errorf("narrate `%s` not defined in config", key)
		}
		narSeq[i] = val
	}

	// ui.Debug("almost")

	srv.player.Speak(&narSeq, srv.config)

	return nil
}

func changeState(state string, srv *ServerState) error {
	logger.Debug("Received state:", state)
	if srv.state == state {
		logger.Debug("already ", state)
		return nil
	}
	srv.state = state
	music, err := srv.config.GetRandomStateMusic(state)
	if err != nil {
		logger.Error(err)
		return err
	}
	srv.player.PlayMusic(music, srv.config, false)
	playlist, err := srv.config.GetStatePlaylist(state)
	if err != nil {
		logger.Error(err)
		return err
	}
	srv.player.SetPlaylist(playlist)

	return nil
}

func changeMusic(state string, srv *ServerState) error {
	logger.Debug("changing music")
	music, err := srv.config.GetRandomStateMusic(state)
	if err != nil {
		logger.Error(err)
		return err
	}
	st, err := srv.config.GetState(state)
	if err != nil {
		logger.Error(err)
		return err
	}
	allowSame := len(st.Music) == 1
	srv.player.PlayMusic(music, srv.config, allowSame)

	return err
}
