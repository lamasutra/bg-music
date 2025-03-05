package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/config"
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
	conf, err := config.Read(*cmdArgs.config)
	if err != nil {
		panic(err)
	}

	initUI(cmdArgs)

	// initServer
	for {
		time.Sleep(time.Millisecond * 100)
		if ui.Ready {
			break
		}
	}

	ui.Debug("Running as", conf.PlayerType, conf.ServerType)
	// debug
	// jsonPretty, _ := json.MarshalIndent(*conf, "", "  ")
	// fmt.Println(string(jsonPretty))

	server, err := server.CreateServer(conf.ServerType)
	if err != nil {
		panic(err)
	}

	defer server.Close()

	server.Serve(conf)
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
