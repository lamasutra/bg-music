package input

import "time"

const EV_INPUT_PRESSED = "input_pressed"

type InputDevice struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type InputManager interface {
	GetInputDevices() ([]InputDevice, error)
	WatchEvents(device InputDevice)
	Close()
	IsWatching() bool
}

type InputManagerWrapper struct {
	manager InputManager
}

func NewManager(debounce time.Duration) *InputManagerWrapper {
	// detect os
	return &InputManagerWrapper{manager: &linuxManager{
		debounce: debounce,
	}}
}

func (m *InputManagerWrapper) GetInputDevices() ([]InputDevice, error) {
	return m.manager.GetInputDevices()
}

func (m *InputManagerWrapper) WatchEvents(device InputDevice) {
	go m.manager.WatchEvents(device)
}

func (m *InputManagerWrapper) IsWatching() bool {
	return m.manager.IsWatching()
}

func (m *InputManagerWrapper) Close() {
	m.manager.Close()
}
