package input

import (
	"time"

	"github.com/lamasutra/bg-music/pkg/logger"
	"github.com/lamasutra/bg-music/wt-client/internal/client"
)

// @todo paralel load
func (p *DataParser) loadData(host string) {
	// fmt.Println("Loading data from", host)
	err := p.inputData.State.Load(host)
	if err != nil {
		// fmt.Println("state error: ", err)
		time.Sleep(sleepOffline)
		err = p.inputData.State.Load(host)
		if err != nil {
			return
		}
	}

	err = p.inputData.Indicators.Load(host)
	if err != nil {
		logger.Error("indicators error: ", err)
	}
	err = p.inputData.MapInfo.Load(host)
	if err != nil {
		logger.Error("!!! mapInfo error: ", err)
	}
	// load other data
	if p.inputData.MapInfo.Valid {
		// load map identity
		if p.inputData.Identity == 0 {
			p.inputData.Identity, err = client.MapIdentity(host)
			logger.Error("map identity error: ", err)
		}

		err = p.inputData.MapObj.Load(host)
		if err != nil {
			logger.Error("mapObj error: ", err)
			// } else {
			// fmt.Println(mapObj)
		}
		err = p.inputData.HudMsg.Load(host, state.lastEvt, state.lastDmg)
		if err != nil {
			logger.Error("hudMsg error: ", err)
		} else {

		}
	} else {
		if p.inputData.Identity != 0 {
			p.inputData.Identity = 0
		}
	}
}
