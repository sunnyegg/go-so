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

var connectedWsClients = make(map[string][]*websocket.Conn)

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

	connectedWsClients[uri.ID] = append(connectedWsClients[uri.ID], ws)

	go chatter(connectedWsClients)
	go eventsub(connectedWsClients)
}

func chatter(wsClients map[string][]*websocket.Conn) {
	chWs := channel.NewChannel(channel.ChannelWebsocket)

	for {
		msgWs := <-chWs.Listen()

		if _, ok := wsClients[msgWs["channel"]]; !ok {
			return
		}

		messageData := WsMessageData{
			Channel:   msgWs["channel"],
			UserLogin: msgWs["username"],
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

		var closedClients []int
		for i, conn := range wsClients[msgWs["channel"]] {
			err = conn.WriteMessage(websocket.TextMessage, msgOutputBytes)
			if err != nil {
				closedClients = append(closedClients, i)
			}
		}

		// remove closed clients
		for _, i := range closedClients {
			wsClients[msgWs["channel"]][i] = wsClients[msgWs["channel"]][len(wsClients[msgWs["channel"]])-1]
			wsClients[msgWs["channel"]] = wsClients[msgWs["channel"]][:len(wsClients[msgWs["channel"]])-1]
		}
	}
}

func eventsub(wsClients map[string][]*websocket.Conn) {
	chEs := channel.NewChannel(channel.ChannelEventsub)

	for {
		msgEs := <-chEs.Listen()

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

		var closedClients []int
		for i, conn := range wsClients[msgEs["channel"]] {
			err = conn.WriteMessage(websocket.TextMessage, msgOutputBytes)
			if err != nil {
				closedClients = append(closedClients, i)
			}
		}

		// remove closed clients
		for _, i := range closedClients {
			wsClients[msgEs["channel"]][i] = wsClients[msgEs["channel"]][len(wsClients[msgEs["channel"]])-1]
			wsClients[msgEs["channel"]] = wsClients[msgEs["channel"]][:len(wsClients[msgEs["channel"]])-1]
		}
	}
}
