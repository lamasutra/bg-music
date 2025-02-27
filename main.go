package main

import (
	"flag"
	"fmt"

	"github.com/lamasutra/bg-music/config"
	"github.com/lamasutra/bg-music/server"
)

// type StateMachine struct {
// 	currentState model.State
// }

type cmdArgs struct {
	config *string
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

	fmt.Println("Running as", conf.PlayerType, conf.ServerType)
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

func registerFlags() *cmdArgs {
	var args cmdArgs
	args.config = flag.String("config", "config.json", "Config file path")

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
