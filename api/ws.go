package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sunnyegg/go-so/channel"
)

type WsMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type WsMessageData struct {
	Channel         string `json:"channel"`
	UserName        string `json:"user_name"`
	UserLogin       string `json:"user_login"`
	Followers       int    `json:"followers"`
	ProfileImageUrl string `json:"profile_image_url"`
	LastSeenPlaying string `json:"last_seen_playing"`
}

type WsURI struct {
	ID string `uri:"id" binding:"required"`
}

var connectedWsClients = make(map[string]map[string]*websocket.Conn)
var chWs = channel.NewChannel(channel.ChannelWebsocket)
var chEs = channel.NewChannel(channel.ChannelEventsub)

func (server *Server) ws(ctx *gin.Context) {
	var uri WsURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	connectedWsClients[uri.ID][ws.RemoteAddr().String()] = ws
}

func (server *Server) listenChannels() {
	go chatter(connectedWsClients)
	go eventsub(connectedWsClients)
}

func chatter(wsClients map[string]map[string]*websocket.Conn) {
	for {
		msgWs := <-chWs.Listen()
		channel := msgWs["channel"]
		userlogin := msgWs["username"]

		if _, ok := wsClients[channel]; !ok {
			return
		}

		messageData := WsMessageData{
			Channel:   channel,
			UserLogin: userlogin,
		}

		// map[string]string to []byte
		msgBytes, err := json.Marshal(messageData)
		if err != nil {
			return
		}

		msgOutput := WsMessage{
			Type: "chatter",
			Data: string(msgBytes),
		}

		msgOutputBytes, err := json.Marshal(msgOutput)
		if err != nil {
			return
		}

		for _, conn := range wsClients[channel] {
			err = conn.WriteMessage(websocket.TextMessage, msgOutputBytes)
			if err != nil {
				conn.Close()
				deleteClient(wsClients, channel, conn.RemoteAddr().String())
			}
		}
	}
}

func eventsub(wsClients map[string]map[string]*websocket.Conn) {
	for {
		msgEs := <-chEs.Listen()
		channel := msgEs["channel"]

		if _, ok := wsClients[channel]; !ok {
			return
		}

		msgBytes, err := json.Marshal(msgEs)
		if err != nil {
			return
		}

		msgOutput := WsMessage{
			Type: "eventsub",
			Data: string(msgBytes),
		}

		msgOutputBytes, err := json.Marshal(msgOutput)
		if err != nil {
			return
		}

		for _, conn := range wsClients[channel] {
			err = conn.WriteMessage(websocket.TextMessage, msgOutputBytes)
			if err != nil {
				conn.Close()
				deleteClient(wsClients, channel, conn.RemoteAddr().String())
			}
		}
	}
}

func deleteClient(wsClients map[string]map[string]*websocket.Conn, channel string, addr string) {
	delete(wsClients[channel], addr)
}
