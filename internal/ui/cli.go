package ui

import (
	"fmt"

	"github.com/lamasutra/bg-music/internal/app"
	"github.com/lamasutra/bg-music/pkg/events"
	"github.com/lamasutra/bg-music/pkg/logger"
)

type cliState struct {
	app *app.AppState
}

func NewCli(a *app.AppState) *cliState {
	return &cliState{
		app: a,
	}
}

func (s *cliState) Run(onStartup func(a *app.AppState)) {
	events.Listen("log", "cli", s.renderMessage)

	onStartup(s.app)
}

func (s *cliState) renderMessage(args ...any) {
	if len(args) != 1 {
		panic("invalid arguments count, cannot render message")
	}

	msg, ok := args[0].(logger.MessageRenderer)
	if !ok {
		panic("message is not renderer")
	}

	fmt.Println(msg.Render())
}
