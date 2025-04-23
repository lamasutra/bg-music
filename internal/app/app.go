package app

import (
	"embed"
	"time"

	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/internal/devices"
	"github.com/lamasutra/bg-music/pkg/input"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

type CmdArgs struct {
	Config  *string
	Tui     *bool
	Cli     *bool
	Verbose *bool
}

type AppState struct {
	Config       *model.Config
	Args         *CmdArgs
	Assets       *embed.FS
	Icon         []byte
	InputManager input.InputManager
	AudioPlayer  audio.Player
}

func NewApp(args *CmdArgs, config *model.Config, assets embed.FS, icon []byte) *AppState {
	return &AppState{
		Config: config,
		Args:   args,
		Assets: &assets,
		Icon:   icon,
	}
}

func (a *AppState) SetInputManager(manager input.InputManager) {
	if a.InputManager != nil {
		a.InputManager.Close()
	}
	a.InputManager = manager
}

func (a *AppState) SetAudioPlayer(ap audio.Player) {
	a.AudioPlayer = ap
}

func (a *AppState) SaveInputControls(payload model.InputDeviceControls) {
	a.Config.SetInputDeviceControls(payload)
	appConfig := &model.Config{}
	appConfig.Read(*a.Args.Config)
	appConfig.SetInputDeviceControls(payload)
	appConfig.Save(*a.Args.Config)
	a.closeKeyboardListener()
	a.RunKeyboardListener()
}

func (a *AppState) closeKeyboardListener() {
	if a.InputManager != nil {
		a.InputManager.Close()
		devices.UnregisterListener()
	}
}

func (a *AppState) RunKeyboardListener() {
	controls := a.Config.Controls
	if controls.Device.Path == "" {
		logger.Warn("no input device specified for keyboard listener")
		return
	}

	inputManager := input.NewManager(200 * time.Millisecond)
	inputManager.WatchEvents(controls.Device.ToManagableDevice())

	devices.ControlPlayer(controls.Mapping, a.AudioPlayer)

	a.SetInputManager(inputManager)
}
