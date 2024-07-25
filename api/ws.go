package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sunnyegg/go-so/channel"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

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
	ch := channel.NewChannel(channel.ChannelWebsocket)

	go func() {
		for {
			msg := <-ch.Listen()
			// map[string]string to []byte
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				return
			}
			err = ws.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				return
			}

			// attendance
			streamID := msg["stream_id"]
			username := msg["username"]
			parsedStreamID, _ := util.ParseStringToInt64(streamID)
			server.store.CreateAttendanceMember(ctx, db.CreateAttendanceMemberParams{
				StreamID:  parsedStreamID,
				Username:  username,
				IsShouted: false,
				PresentAt: util.StringToTimestamp(time.Now().Format(time.RFC3339)),
			})
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
