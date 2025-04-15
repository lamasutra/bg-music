package gui

import (
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/pkg/events"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

type PlayerControls struct {
	music    model.Music
	metadata model.MusicMetadata
}

func NewPlayerControls() *PlayerControls {
	pc := PlayerControls{}
	logger.Debug("Listening music event")
	events.Listen(audio.EV_MUSIC, "player_controls", func(args ...any) { pc.onMusicChange(args[0], args[1]) })
	return &pc
}

func (p *PlayerControls) SetVolume(volume uint8) {
	events.Trigger(audio.EV_SET_VOLUME, volume)
}

func (p *PlayerControls) Next() {
	events.Trigger(audio.EV_PLAY_NEXT)
}

func (p *PlayerControls) Prev() {
	events.Trigger(audio.EV_PLAY_PREVIOUS)
}

func (p *PlayerControls) Pause() {
	events.Trigger(audio.EV_TOGGLE_PAUSE)
}

func (p *PlayerControls) Mute() {
	events.Trigger(audio.EV_TOGGLE_MUTE)
}

func (p *PlayerControls) onMusicChange(musicAny any, metadataAny any) {
	fmt.Println("music change", musicAny, metadataAny, time.Now().UnixMilli()/1000)
	music, ok := musicAny.(model.Music)
	if !ok {
		logger.Debug("music not ok", metadataAny)
		return
	}
	metadata, ok := metadataAny.(model.MusicMetadata)
	if !ok {
		logger.Debug("metadata not ok", metadataAny)
		return
	}
	p.music = music
	p.metadata = metadata
}

func (p *PlayerControls) GetCurrentMusic() model.Music {
	return p.music
}

func (p *PlayerControls) GetCurrentMetadata() model.MusicMetadata {
	return p.metadata
}
