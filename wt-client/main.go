package main

import (
	"fmt"
	"time"

	"github.com/lamasutra/bg-music/wt-client/input"
	"github.com/lamasutra/bg-music/wt-client/model"
	"github.com/lamasutra/bg-music/wt-client/player"
	"github.com/lamasutra/bg-music/wt-client/stateMachine"
)

func main() {
	var conf model.Config
	err := conf.Read("wt-config.json")
	if err != nil {
		fmt.Println("Cannot read wt-config.json")
		return
	}

	fmt.Println("your nickname", conf.Nickname)

	bgPayer := player.CreatePlayer(conf.BgPlayerType, &conf)
	stMachine := stateMachine.New("idle", &conf.StateRules)

	// debug
	fmt.Println(stMachine)

	defer bgPayer.Close()
	for {
		input.LoadLoop(conf.Host, &conf, stMachine, bgPayer)
		time.Sleep(time.Millisecond * 100)
		// fmt.Println("tick")
	}
}
