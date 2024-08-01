package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sunnyegg/go-so/channel"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

type WsMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (server *Server) ws(ctx *gin.Context) {
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

	err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		fmt.Println("error")
		return
	}

	// connect to channel
	chWs := channel.NewChannel(channel.ChannelWebsocket)
	chEs := channel.NewChannel(channel.ChannelEventsub)

	go func() {
		for {
			msgWs := <-chWs.Listen()

			// attendance
			streamID := msgWs["stream_id"]
			username := msgWs["username"]
			parsedStreamID, _ := util.ParseStringToInt64(streamID)
			_, err = server.store.CreateAttendanceMember(ctx, db.CreateAttendanceMemberParams{
				StreamID:  parsedStreamID,
				Username:  username,
				IsShouted: false,
				PresentAt: util.StringToTimestamp(time.Now().Format(time.RFC3339)),
			})
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					fmt.Println("member exists")
					return
				}

				fmt.Println(err)
			}

			// map[string]string to []byte
			msgBytes, err := json.Marshal(msgWs)
			if err != nil {
				return
			}

			msgOutput := WsMessage{
				Type: "attendance",
				Data: string(msgBytes),
			}

			msgOutputBytes, err := json.Marshal(msgOutput)
			if err != nil {
				return
			}

			err = ws.WriteMessage(websocket.TextMessage, msgOutputBytes)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		for {
			msgEs := <-chEs.Listen()

			msgBytes, err := json.Marshal(msgEs)
			if err != nil {
				return
			}
			err = ws.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				return
			}
		}
	}()

	reader(ws)
}

func reader(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("Received message: %s\n", message)
	}
}
