package server

import (
	"encoding/json"
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

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = *ui.GetState()

	router := gin.Default()
	router.POST("/control/:action", controlHandler)
	router.PUT("/state/:code", stateHandler)
	router.PUT("/event/:code", eventHandler)

	go polishAladinsLamp(router)

	for {
		select {
		case <-(*serverState.player).GetMusicEndedChan():
			changeMusic(h.state.state, h.state)
		default:
			time.Sleep(sleepTime)
		}
	}
}

func polishAladinsLamp(router *gin.Engine) {
	err := router.Run(":8211")
	if err != nil {
		panic(err)
	}
}

func (h *HttpServer) loadConfig(data *LoadData) {
	h.state.config.Events = data.Events
	h.state.config.States = data.States

	str, _ := json.MarshalIndent(h.state.config, "", "  ")
	ui.Debug(string(str))
}

func controlHandler(c *gin.Context) {
	action := c.Param("action")

	ui.Debug("Received control:", action)

	switch action {
	case "load":
		data := LoadData{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		instance.loadConfig(&data)
		c.Status(http.StatusNoContent)
		return
	case "next":
		changeMusic(instance.state.state, instance.state)
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
}

func stateHandler(c *gin.Context) {
	state := c.Param("code")
	changeState(state, instance.state)
	c.Status(http.StatusNoContent)
}

func eventHandler(c *gin.Context) {
	event := c.Param("code")
	triggerEvent(event, instance.state)
	c.Status(http.StatusNoContent)
}

func (h *HttpServer) Close() {

}
