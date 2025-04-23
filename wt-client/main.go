package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"flag"
	"fmt"

	"github.com/lamasutra/bg-music/pkg/events"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/wt-client/internal/input"
	"github.com/lamasutra/bg-music/wt-client/internal/model"
	"github.com/lamasutra/bg-music/wt-client/internal/player"
	"github.com/lamasutra/bg-music/wt-client/internal/ui"
)

type cmdArgs struct {
	tui *bool
}

func main() {
	logger.SetLevel(logger.LevelTrace)
	events.Listen("log", "cli", renderMessage)

	go func() {
		log.Println(http.ListenAndServe("localhost:6062", nil))
	}()

	cmdArgs := registerFlags()
	if cmdArgs == nil {
		return
	}

	initUI(cmdArgs)

	var conf model.Config
	err := conf.Read("wt-config.json")
	if err != nil {
		logger.Error("Cannot read wt-config.json")
		return
	}
	err = conf.StateRules.Read("rules.json")
	if err != nil {
		logger.Error("Cannot read rules.json")
		return
	}

	logger.Debug("your configured nickname", "`"+conf.Nickname+"`")

	bgPayer := player.CreatePlayer(conf.BgPlayerType, &conf)
	stMachine := model.NewStateMachine("idle", &conf.StateRules)
	inputLoop := input.CreateInputLoop(&conf, stMachine, bgPayer)

	// debug
	logger.Debug(stMachine)

	defer bgPayer.Close()
	inputLoop.Run()
}

func renderMessage(args ...any) {
	if len(args) != 1 {
		panic("invalid arguments count, cannot render message")
	}

	msg, ok := args[0].(logger.MessageRenderer)
	if !ok {
		panic("message is not renderer")
	}

	fmt.Println(msg.Render())
}

func initUI(args *cmdArgs) {
	uiType := "cli"
	if *args.tui {
		uiType = "tui"
	}
	ui.CreateUI(uiType)
}

func registerFlags() *cmdArgs {
	var args cmdArgs
	args.tui = flag.Bool("tui", false, "show tui")

	// Use a flag with usage function as its value
	helpFlag := flag.Bool("h", false, usage())
	versionFlag := flag.Bool("v", false, "")
	flag.Parse()

	if *helpFlag {
		fmt.Println(usage())
		return nil
	} else if *versionFlag {
		fmt.Println("version: poc")
		return nil
	}

	return &args
}

func usage() string {
	return `
Usage:
  -h|--help   Show this message and exit
  -v          Print version information
  --tui       Render text user interface

Flags:
  config	The config file path (defauklt: "config.json")
`
}
