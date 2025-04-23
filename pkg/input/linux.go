package input

import (
	"path/filepath"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/lamasutra/bg-music/pkg/events"
	"github.com/lamasutra/bg-music/pkg/logger"
)

type linuxManager struct {
	watching, terminate bool
	debounce            time.Duration
	fnKeys              struct {
		ctrl  bool
		alt   bool
		shift bool
	}
	lastPressedKey evdev.EvCode
	lastKey        string
	lastKeyTsms    int64
}

func (l *linuxManager) GetInputDevices() ([]InputDevice, error) {
	var devices []InputDevice

	matches, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return nil, err
	}

	for _, devPath := range matches {
		dev, err := evdev.Open(devPath)
		if err != nil {
			continue
		}

		name, err := dev.Name()
		if err != nil {
			continue
		}

		devices = append(devices, InputDevice{
			Name: name,
			Path: devPath,
		})
	}

	return devices, nil
}

func (l *linuxManager) WatchEvents(device InputDevice) {
	sleepTime := time.Millisecond * 50
	l.terminate = false
	l.watching = true
	dev, err := evdev.Open(device.Path)
	// fmt.Println("WatchEvents on " + device.Path)
	if err != nil {
		// fmt.Println(fmt.Sprintf("Failed to open device %s", device.Path))
		logger.Fatal(err)
		l.watching = false
		return
	}

	defer dev.Close()
	for {
		if l.terminate {
			break
		}

		event, err := dev.ReadOne()

		// fmt.Println(event)

		if err != nil {
			// fmt.Println(fmt.Sprintf("Failed to read from device %s", device.Path))
			logger.Error("Read error: %v", err)
			continue
		}

		// fmt.Println(event)

		if event.Type != evdev.EV_KEY {
			continue
		}

		switch event.Value {
		case 1: // key down
			switch event.Code {
			case evdev.KEY_LEFTCTRL, evdev.KEY_RIGHTCTRL:
				l.fnKeys.ctrl = true
			case evdev.KEY_LEFTALT, evdev.KEY_RIGHTALT:
				l.fnKeys.alt = true
			case evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT:
				l.fnKeys.shift = true
			default:
				l.lastPressedKey = event.Code
			}
		case 0: // key up
			switch event.Code {
			case evdev.KEY_LEFTCTRL, evdev.KEY_RIGHTCTRL:
				l.fnKeys.ctrl = false
			case evdev.KEY_LEFTALT, evdev.KEY_RIGHTALT:
				l.fnKeys.alt = false
			case evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT:
				l.fnKeys.shift = false
			default:
				l.lastPressedKey = evdev.KEY_RESERVED
			}
		}
		l.triggerKeyChange()
		time.Sleep(sleepTime)
	}
	l.watching = false
}

func (l *linuxManager) IsWatching() bool {
	return l.watching
}

func (l *linuxManager) Close() {
	l.terminate = true
}

func (l *linuxManager) triggerKeyChange() {
	combo := ""
	if l.fnKeys.ctrl {
		combo += "CTRL+"
	}
	if l.fnKeys.alt {
		combo += "ALT+"
	}
	if l.fnKeys.shift {
		combo += "SHIFT+"
	}
	if l.lastPressedKey != evdev.KEY_RESERVED {
		key := evdev.CodeName(evdev.EV_KEY, l.lastPressedKey)
		combo += key[4:]
	} else if len(combo) > 0 {
		// get rid of +
		combo = combo[0 : len(combo)-1]
	}

	// fmt.Println("C:" + combo)
	nowTsms := time.Now().UnixMilli()
	if combo != l.lastKey {
		l.lastKey = combo
		logger.Trace("Triggering key change: " + combo)
		l.lastKeyTsms = nowTsms
		events.Trigger(EV_INPUT_PRESSED, combo)
	} else {
		if l.lastKeyTsms != 0 && (nowTsms-l.lastKeyTsms) > int64(l.debounce) {
			logger.Trace("Triggering key press after debounce: " + combo)
			events.Trigger(EV_INPUT_PRESSED, combo)
		}
	}
}
