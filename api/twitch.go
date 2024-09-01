package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sunnyegg/go-so/channel"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

var connectedTwitchClients = make(map[string]bool)

func (server *Server) getTwitchUser(ctx *gin.Context) {
	var req getTwitchUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, _, err := decryptHeader(ctx, server)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)
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

	// skip if already connected
	if _, ok := connectedTwitchClients[req.UserLogin]; ok {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	payload, authPayload, err := decryptHeader(ctx, server)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// connect to chat
	configChat := twitch.ConnectConfig{
		StreamID: req.StreamID,
		IsAutoSO: false,
		Delay:    5,
	}

	twChatClient := twitch.NewChatClient(req.UserLogin, payload.AccessToken)
	twChatClient.Connect(configChat)
	twChatClient.Join(req.UserLogin, req.Channel)

	connectedTwitchClients[req.UserLogin] = true

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

	ch := channel.NewChannel(channel.ChannelEventsub)

	// if someone goes live
	if req.Subscription.Type == EventsubSubscriptionTypeStreamOnline {
		go ch.Send(map[string]string{
			"type":    EventsubSubscriptionTypeStreamOnline,
			"channel": req.Event.UserLogin,
		})

		ctx.JSON(http.StatusOK, nil)
		return
	}

	// if someone goes offline
	if req.Subscription.Type == EventsubSubscriptionTypeStreamOffline {
		delete(connectedTwitchClients, req.Event.UserLogin)
		go ch.Send(map[string]string{
			"type":    EventsubSubscriptionTypeStreamOffline,
			"channel": req.Event.UserLogin,
		})

		ctx.JSON(http.StatusOK, nil)
		return
	}

	// if someone redeems reward
	if req.Subscription.Type == EventsubSubscriptionTypeChannelRedemption {
		go ch.Send(map[string]string{
			"type":     EventsubSubscriptionTypeChannelRedemption,
			"channel":  req.Event.BroadcasterUserLogin,
			"redeemer": req.Event.UserLogin,
		})

		ctx.JSON(http.StatusOK, nil)
		return
	}
}

func (server *Server) getChannelInfo(ctx *gin.Context) {
	var req getChannelInfoRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, _, err := decryptHeader(ctx, server)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)
	userInfo, err := twClient.GetUserInfo(payload.AccessToken, "", req.UserLogin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelInfo, err := twClient.GetChannelInfo(payload.AccessToken, userInfo.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelFollowers, err := twClient.GetChannelFollowers(payload.AccessToken, userInfo.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelInfo.Followers = channelFollowers.Total
	channelInfo.User.Login = userInfo.Login
	channelInfo.User.DisplayName = userInfo.DisplayName
	channelInfo.User.ProfileImageURL = userInfo.ProfileImageURL

	ctx.JSON(http.StatusOK, channelInfo)
}

func (server *Server) getStreamInfo(ctx *gin.Context) {
	var req getChannelInfoRequest
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

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)
	userInfo, err := twClient.GetUserInfo(payload.AccessToken, "", req.UserLogin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	streamInfo, err := twClient.GetStreamInfo(payload.AccessToken, userInfo.ID)
	if err != nil {
		if err.Error() == "stream not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, streamInfo)
}

func (server *Server) sendChatMessage(ctx *gin.Context) {
	var req sendChatMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, _, err := decryptHeader(ctx, server)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)

	// get user id
	channelInfo, err := twClient.GetUserInfo(payload.AccessToken, "", req.Channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient.SendChatMessage(payload.AccessToken, channelInfo.ID, channelInfo.ID, req.Message)

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) sendShoutout(ctx *gin.Context) {
	var req sendShoutoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, _, err := decryptHeader(ctx, server)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)

	// get user id
	fromUser, err := twClient.GetUserInfo(payload.AccessToken, "", req.FromID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	toUser, err := twClient.GetUserInfo(payload.AccessToken, "", req.ToID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	moderatorUser, err := twClient.GetUserInfo(payload.AccessToken, "", req.ModeratorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	twClient.SendShoutout(payload.AccessToken, fromUser.ID, toUser.ID, moderatorUser.ID)

	ctx.JSON(http.StatusOK, nil)
}
