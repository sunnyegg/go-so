package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sunnyegg/go-so/channel"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

var ConnectedClients = make(map[string]bool)

func (server *Server) getTwitchUser(ctx *gin.Context) {
	var req getTwitchUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get sessionid
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// session
	session, err := server.store.GetSession(ctx, db.GetSessionParams{
		ID:     util.UUIDToUUID(authPayload.SessionID),
		UserID: authPayload.UserID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// decrypt token
	tokenBytes, err := util.Decrypt(session.EncryptedTwitchToken, server.config.TokenSymmetricKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token := []byte(tokenBytes)
	payload := twitch.OAuthToken{}
	err = json.Unmarshal(token, &payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.RedirectURI)
	userInfo, err := twClient.GetUserInfo(payload.AccessToken, "", req.UserLogin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

func (server *Server) connectChat(ctx *gin.Context) {
	var req connectChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get sessionid
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// skip if already connected
	if _, ok := ConnectedClients[req.UserLogin]; ok {
		return
	}

	// session
	session, err := server.store.GetSession(ctx, db.GetSessionParams{
		ID:     util.UUIDToUUID(authPayload.SessionID),
		UserID: authPayload.UserID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// decrypt token
	tokenBytes, err := util.Decrypt(session.EncryptedTwitchToken, server.config.TokenSymmetricKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token := []byte(tokenBytes)
	payload := twitch.OAuthToken{}
	err = json.Unmarshal(token, &payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// connect to chat
	configChat := twitch.ConnectConfig{
		StreamID: req.StreamID,
		IsAutoSO: false,
		Delay:    5,
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.RedirectURI)
	twClient.ConnectTwitchChat(configChat, req.UserLogin, req.UserLogin, payload.AccessToken)

	ConnectedClients[req.UserLogin] = true

	// get user config blacklist
	userConfig, err := server.store.GetUserConfig(ctx, db.GetUserConfigParams{
		UserID:     authPayload.UserID,
		ConfigType: db.ConfigTypesBlacklist,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusOK, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// send list blacklist to twitch chat
	ch := channel.NewChannel(channel.ChannelBlacklist)
	ch.Send(map[string]string{
		req.UserLogin: userConfig.Value,
	})

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) handleEventsub(ctx *gin.Context) {
	var req eventsubRequest
	var err error

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: validate signature
	if ctx.Request.Header.Get(EventsubMessageTypeHeaderKey) != "notification" {
		if ctx.Request.Header.Get(EventsubMessageTypeHeaderKey) == "webhook_callback_verification" {
			ctx.Request.Header.Set("Content-Type", "text/plain")
			ctx.String(http.StatusOK, req.Challenge)
			return
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid message type")))
		return
	}

	// get user by userid
	user, err := server.store.GetUserByUserID(ctx, req.Event.BroadcasterUserID)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, nil)
		return
	}

	// if someone goes live
	if req.Subscription.Type == EventsubSubscriptionTypeStreamOnline {
		ch := channel.NewChannel(channel.ChannelEventsub)
		go func() {
			ch.Send(map[string]string{
				req.Subscription.Condition.BroadcasterUserID: "stream online",
			})
		}()

		ctx.JSON(http.StatusOK, nil)
		return
	}

	// if someone goes offline
	if req.Subscription.Type == EventsubSubscriptionTypeStreamOffline {
		delete(ConnectedClients, req.Event.UserLogin)

		ctx.JSON(http.StatusOK, nil)
		return
	}

	// if someone redeems reward
	if req.Subscription.Type == EventsubSubscriptionTypeChannelRedemption {
		streams, err := server.store.ListStreams(ctx, db.ListStreamsParams{
			UserID: user.ID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ch := channel.NewChannel(channel.ChannelWebsocket)
		go func() {
			ch.Send(map[string]string{
				"stream_id": util.ParseIntToString(int(streams[0].ID)),
				"username":  req.Event.UserLogin,
			})
		}()

		ctx.JSON(http.StatusOK, nil)
		return
	}
}
