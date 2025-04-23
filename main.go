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
	"github.com/lamasutra/bg-music/internal/app"
	"github.com/lamasutra/bg-music/internal/audio"
	"github.com/lamasutra/bg-music/internal/ui"
	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/pkg/model"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/icons/music-app-systray-icon.ico
var icon []byte

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	cmdArgs := registerFlags()
	if cmdArgs == nil {
		return
	}
	config := &model.Config{}
	err := config.Read(*cmdArgs.Config)
	if err != nil {
		panic(err)
	}

	a := app.NewApp(cmdArgs, config, assets, icon)

	createUI(a, func(app *app.AppState) {
		onStartup(app)
	})
}

func onStartup(a *app.AppState) {
	time.Sleep(time.Second)

	mp := audio.CreatePlayer(a.Config.PlayerType)

	defer mp.Close()

	go runServer(a.Config, mp)

	a.SetAudioPlayer(mp)
	go a.RunKeyboardListener()

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

func createUI(a *app.AppState, onStartup func(a *app.AppState)) {
	uiType := "gui"
	if *a.Args.Tui {
		uiType = "tui"
	} else if *a.Args.Cli {
		uiType = "cli"
	}

	if *a.Args.Verbose {
		logger.SetLevel(logger.LevelInfo)
	} else {
		logger.SetLevel(logger.LevelWarn)
	}

	ui.CreateNew(uiType, a, onStartup)

}

func registerFlags() *app.CmdArgs {
	var args app.CmdArgs
	args.Config = flag.String("config", "config.json", "Config file path")
	args.Tui = flag.Bool("tui", false, "show tui")
	args.Cli = flag.Bool("cli", false, "pure cli")
	args.Verbose = flag.Bool("vv", false, "verbose mode")

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
  -vv         Verbose mode, print more debug info

Flags:
  config	The config file path (defauklt: "config.json")
`
}
