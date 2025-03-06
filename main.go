package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/model"
	"github.com/lamasutra/bg-music/player"
	"github.com/lamasutra/bg-music/server"
	"github.com/lamasutra/bg-music/ui"
)

// type StateMachine struct {
// 	currentState model.State
// }

type cmdArgs struct {
	config *string
	tui    *bool
}

func main() {
	cmdArgs := registerFlags()
	if cmdArgs == nil {
		return
	}
	config := &model.Config{}
	err := config.Read(*cmdArgs.config)
	if err != nil {
		panic(err)
	}

	go initUI(cmdArgs)

	time.Sleep(time.Second)

	mp := player.CreatePlayer(config.PlayerType)

	defer (*mp).Close()

	go initServer(config, mp)

	for {
		time.Sleep(time.Second)
	}
}

func initServer(config *model.Config, mp *model.Player) {
	ui.Debug("Running as", config.PlayerType, config.ServerType)

	server, err := server.CreateServer(config.ServerType)
	if err != nil {
		panic(err)
	}

	defer server.Close()

	server.Serve(config, mp)
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
	args.config = flag.String("config", "config.json", "Config file path")
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

Flags:
  config	The config file path (defauklt: "config.json")
`
}
