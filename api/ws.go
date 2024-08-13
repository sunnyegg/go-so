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
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
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

var connectedClients = make(map[string][]*websocket.Conn)

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

	connectedClients[uri.ID] = append(connectedClients[uri.ID], ws)

	// connect to channel
	chWs := channel.NewChannel(channel.ChannelWebsocket)
	chEs := channel.NewChannel(channel.ChannelEventsub)

	go func() {
		for {
			msgWs := <-chWs.Listen()

			if _, ok := connectedClients[msgWs["channel"]]; !ok {
				return
			}

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

			twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)
			userInfo, err := twClient.GetUserInfo(msgWs["token"], "", msgWs["username"])
			if err != nil {
				fmt.Println(err)
				return
			}

			channelInfo, err := twClient.GetChannelInfo(msgWs["token"], userInfo.ID)
			if err != nil {
				fmt.Println(err)
				return
			}

			channelFollowers, err := twClient.GetChannelFollowers(msgWs["token"], userInfo.ID)
			if err != nil {
				fmt.Println(err)
				return
			}

			messageData := WsMessageData{
				Channel:         msgWs["channel"],
				UserName:        userInfo.DisplayName,
				UserLogin:       userInfo.Login,
				Followers:       channelFollowers.Total,
				ProfileImageUrl: userInfo.ProfileImageURL,
				LastSeenPlaying: channelInfo.GameName,
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

			for _, conn := range connectedClients[msgWs["channel"]] {
				err = conn.WriteMessage(websocket.TextMessage, msgOutputBytes)

				// TODO: delete client if error
				if err != nil {
					conn.Close()
				}
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
	var msg WsMessage

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(msg.Data)
	}
}
