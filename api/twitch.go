package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

type getTwitchUserRequest struct {
	Username string `form:"username" binding:"required"`
}

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
		ID:     util.UUIDToUUID(authPayload.ID),
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
	userInfo, err := twClient.GetUserInfo(payload.AccessToken, "", req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

type connectChatRequest struct {
	StreamID string `json:"stream_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Channel  string `json:"channel" binding:"required"`
}

func (server *Server) connectChat(ctx *gin.Context) {
	var req connectChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get sessionid
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// session
	session, err := server.store.GetSession(ctx, db.GetSessionParams{
		ID:     util.UUIDToUUID(authPayload.ID),
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
		IsAutoSO: true,
		Delay:    5,
	}

	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.RedirectURI)
	twClient.ConnectTwitchChat(configChat, req.Username, payload.AccessToken)

	ctx.JSON(http.StatusOK, nil)
}

type eventsubRequest struct {
	Challenge    string `json:"challenge"`
	Subscription struct {
		Type      string `json:"type"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		}
	} `json:"subscription"`
	Event struct {
		UserLogin string `json:"user_login"`
		Reward    struct {
			Title string `json:"title"`
		} `json:"reward"`
	} `json:"event"`
}

const (
	EventsubMessageIDHeaderKey                = "Twitch-Eventsub-Message-Id"
	EventsubMessageTimestampHeaderKey         = "Twitch-Eventsub-Message-Timestamp"
	EventsubMessageSignatureHeaderKey         = "Twitch-Eventsub-Message-Signature"
	EventsubMessageTypeHeaderKey              = "Twitch-Eventsub-Message-Type"
	EventsubSubscriptionTypeChannelRedemption = "channel.channel_points_custom_reward_redemption.add"
	EventsubSubscriptionTypeStreamOnline      = "stream.online"
	EventsubSubscriptionTypeFollow            = "channel.follow"
)

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

	// if someone goes live
	// create stream
	if req.Subscription.Type == EventsubSubscriptionTypeStreamOnline {
		// get session by userid
		session, err := server.store.GetSessionByUserID(ctx, req.Subscription.Condition.BroadcasterUserID)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

		// get user stream info
		twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.RedirectURI)
		streamInfo, err := twClient.GetStreamInfo(payload.AccessToken, req.Subscription.Condition.BroadcasterUserID)
		if err != nil {
			// if stream not found, not doing anything
			if err.Error() == "stream not found" {
				ctx.JSON(http.StatusOK, nil)
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// create stream
		arg := db.CreateStreamParams{
			UserID:    session.UserID,
			Title:     streamInfo.Title,
			GameName:  streamInfo.GameName,
			StartedAt: util.StringToTimestamp(streamInfo.StartedAt),
			CreatedBy: session.UserID,
		}

		stream, err := server.store.CreateStream(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// connect chat
		configChat := twitch.ConnectConfig{
			StreamID: util.ParseIntToString(int(stream.ID)),
			IsAutoSO: true,
			Delay:    5,
			Blacklist: []string{
				"NIGHTBOT",
			},
		}

		twClient.ConnectTwitchChat(configChat, req.Event.UserLogin, payload.AccessToken)
	}

	ctx.JSON(http.StatusOK, nil)
}
