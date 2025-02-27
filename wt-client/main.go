package main

import (
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/wt-client/clientConfig"
	"github.com/lamasutra/bg-music/wt-client/input"
	"github.com/lamasutra/bg-music/wt-client/player"
	"github.com/lamasutra/bg-music/wt-client/stateMachine"
)

func main() {
	var conf clientConfig.Config
	err := conf.Read("wt-config.json")
	if err != nil {
		fmt.Println("Cannot read wt-config.json")
		return
	}

	fmt.Println("your nickname", conf.Nickname)

	player := player.CreatePlayer(conf.BgPlayerType, &conf)
	stMachine := stateMachine.New("idle", &conf.StateRules)

	// debug
	// fmt.Println(stMachine)

	defer player.Close()
	for {
		input.LoadLoop(conf.Host, &conf, stMachine, player)
		time.Sleep(time.Millisecond * 50)
		// fmt.Println("tick")
	}
}
