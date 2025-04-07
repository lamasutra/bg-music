package input

import (
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

// @todo paralel load
func loadData(host string) {
	// fmt.Println("Loading data from", host)
	err := inputData.State.Load(host)
	if err != nil {
		// fmt.Println("state error: ", err)
		time.Sleep(sleepOffline)
		err = inputData.State.Load(host)
		if err != nil {
			return
		}
	}

	err = inputData.Indicators.Load(host)
	if err != nil {
		ui.Error("indicators error: ", err)
	}
	err = inputData.MapInfo.Load(host)
	if err != nil {
		ui.Error("!!! mapInfo error: ", err)
	}
	// load other data
	if inputData.MapInfo.Valid {
		// load map identity
		if inputData.Identity == 0 {
			inputData.Identity, err = client.MapIdentity(host)
			ui.Error("map identity error: ", err)
		}

		err = inputData.MapObj.Load(host)
		if err != nil {
			ui.Error("mapObj error: ", err)
			// } else {
			// fmt.Println(mapObj)
		}
		err = inputData.HudMsg.Load(host, state.lastEvt, state.lastDmg)
		if err != nil {
			ui.Error("hudMsg error: ", err)
		} else {

		}
	} else {
		if inputData.Identity != 0 {
			inputData.Identity = 0
		}
	}
}
