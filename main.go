package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/lamasutra/bg-music/internal/api"
	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/internal/devices"
	"github.com/lamasutra/bg-music/internal/ui"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/icons/music-app-systray-icon.ico
var icon []byte

type cmdArgs struct {
	config *string
	tui    *bool
	cli    *bool
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	cmdArgs := registerFlags()
	if cmdArgs == nil {
		return
	}
	config := &model.Config{}
	err := config.Read(*cmdArgs.config)
	if err != nil {
		panic(err)
	}

	createUI(cmdArgs, &assets, icon, func() {
		onStartup(config)
	})
}

func onStartup(config *model.Config) {
	time.Sleep(time.Second)

	mp := audio.CreatePlayer(config.PlayerType)

	defer mp.Close()

	go runServer(config, mp)

	go runKeyboardListener(config.Controls, mp)

	for {
		time.Sleep(time.Second)
	}
}

func runServer(config *model.Config, mp audio.Player) {
	logger.Debug("Running as ", config.PlayerType, " ", config.ServerType)

	server, err := api.CreateServer(config.ServerType)
	if err != nil {
		panic(err)
	}

	defer server.Close()

	server.Serve(config, mp)
}

func createUI(args *cmdArgs, assets *embed.FS, icon []byte, onStartup func()) {
	uiType := "gui"
	if *args.tui {
		uiType = "tui"
	} else if *args.cli {
		uiType = "cli"
	}

	ui.CreateNew(uiType, assets, icon, onStartup)
}

func runKeyboardListener(controls map[string]string, mp audio.Player) {
	devices.WatchInput(controls, mp)
}

func registerFlags() *cmdArgs {
	var args cmdArgs
	args.config = flag.String("config", "config.json", "Config file path")
	args.tui = flag.Bool("tui", false, "show tui")
	args.cli = flag.Bool("cli", false, "pure cli")

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
  --cli       Render GUI interface

Flags:
  config	The config file path (defauklt: "config.json")
`
}
