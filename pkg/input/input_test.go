package input

import (
	"testing"
	"time"

	"github.com/lamasutra/bg-music/pkg/events"
)

func TestInput(t *testing.T) {
	manager := NewManager(200 * time.Millisecond)
	devices, err := manager.GetInputDevices()
	if err != nil {
		t.Errorf("Failed to get input devices: %v", err)
	}
	for _, device := range devices {
		t.Logf("Device Name: %s, Path: %s", device.Name, device.Path)
	}

	device := InputDevice{Name: "Keyboard", Path: "/dev/input/event11"}
	manager.WatchEvents(device)
	events.Listen("input_pressed", "test", func(args ...any) {
		combo, ok := args[0].(string)
		if !ok {
			t.Errorf("Failed to cast input_pressed argument as string")
			return
		}
		t.Logf("Keyboard combo pressed: %s", combo)
	})
	time.Sleep(time.Second * 10)
	manager.Close()
}
