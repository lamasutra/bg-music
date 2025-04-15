package devices

import (
	"bufio"
	"os/exec"
	"time"

	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/pkg/logger"
)

func WatchInput(controls map[string]string, player audio.Player) {
	for {
		pyInputListener(controls, player)
		time.Sleep(time.Second)
	}
}

func pyInputListener(controls map[string]string, player audio.Player) {
	cmd := exec.Command("python3", "listener.py")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	logger.Debug("Listening for global keypresses (Ctrl+C to quit)...")

	for scanner.Scan() {
		combo := scanner.Text()
		logger.Debug("Key Pressed:", combo)
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
	logger.Error("Stopped listening for global keypresses.")
}
