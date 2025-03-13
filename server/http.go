package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lamasutra/bg-music/model"
	"github.com/lamasutra/bg-music/ui"
)

type HttpServer struct {
	state *ServerState
}

var instance *HttpServer

func NewHttpServer() *HttpServer {
	instance = &HttpServer{
		state: &ServerState{},
	}

	return instance
}

func (h *HttpServer) Serve(conf *model.Config, player *model.Player) {
	sleepTime := time.Millisecond * 100
	serverState := ServerState{
		config: conf,
		player: player,
	}
	// player.CreatePlayer(conf.PlayerType, conf.Volume, &p.musicEndedChannel),
	h.state = &serverState

	defer (*h.state.player).Close()

	changeState("idle", h.state)

	router := gin.Default()
	router.POST("/control", controlHandler)
	router.POST("/state", stateHandler)
	router.POST("/event", eventHandler)

	router.Run(":8211")

	for {
		select {
		case <-(*serverState.player).GetMusicEndedChan():
			changeMusic(h.state.state, h.state)
		default:
			time.Sleep(sleepTime)
		}
	}
}

func (h *HttpServer) loadConfig(data *LoadData) {
	h.state.config.Events = data.Events
	h.state.config.States = data.States

	str, _ := json.MarshalIndent(h.state.config, "", "  ")
	ui.Debug(string(str))
}

func controlHandler(c *gin.Context) {
	req := Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ui.Debug("Received control:", req.Action)

	switch req.Action {
	case "load":
		loadRequest := LoadRequest{}
		if err := c.ShouldBindJSON(&loadRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		instance.loadConfig(&loadRequest.Data)
	case "next":
		changeMusic(instance.state.state, instance.state)
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})

}

func stateHandler(c *gin.Context) {
	req := StateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	changeState(req.State, instance.state)
}

func eventHandler(c *gin.Context) {
	req := EventRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(instance)
	triggerEvent(req.Event, instance.state)
}

func (h *HttpServer) Close() {

}
