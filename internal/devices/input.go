package devices

import (
	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/pkg/events"
	"github.com/lamasutra/bg-music/pkg/input"
	"github.com/lamasutra/bg-music/pkg/logger"
)

func ControlPlayer(mapping map[string]string, player audio.Player) {
	logger.Info(mapping)
	events.Listen(input.EV_INPUT_PRESSED, "controls", func(args ...any) {
		logger.Debug(input.EV_INPUT_PRESSED, args)
		combo, ok := args[0].(string)
		if !ok {
			logger.Warn("not a string", args[0])
			return
		}
		ctrl, ok := mapping[combo]
		if !ok {
			return
		}
		player.SendControl(ctrl)
	})
}

func UnregisterListener() {
	events.Stop(input.EV_INPUT_PRESSED, "controls")
}
