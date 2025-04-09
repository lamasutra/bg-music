package devices

import (
	"bufio"
	"os/exec"
	"time"

	"github.com/lamasutra/bg-music/internal/ui"
	"github.com/lamasutra/bg-music/model"
)

func WatchInput(controls map[string]string, player model.Player) {
	for {
		pyInputListener(controls, player)
		time.Sleep(time.Second)
	}
}

func pyInputListener(controls map[string]string, player model.Player) {
	cmd := exec.Command("python3", "listener.py")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	ui.Debug("Listening for global keypresses (Ctrl+C to quit)...")

	for scanner.Scan() {
		combo := scanner.Text()
		ui.Debug("Key Pressed:", combo)
		ctrl, ok := controls[combo]
		if !ok {
			continue // Ignore unknown keys
		}
		player.SendControl(ctrl)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	cmd.Wait()
	ui.Error("Stopped listening for global keypresses.")
}
