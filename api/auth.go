package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

var tempState = make(map[string]bool, 0)

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Error != "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	// check state
	if _, ok := tempState[req.State]; !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid state")))
		return
	}
	delete(tempState, req.State)

	// login twitch
	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.FeAddress)
	token, err := twClient.GetOAuthToken(req.Code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// validate token
	err = twClient.ValidateOAuthToken(token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get user info twitch
	userInfo, err := twClient.GetUserInfo(token.AccessToken, "", "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check user login
	// if not exists, createUser
	// else updateUser
	rsp, err := createOrUpdateUser(ctx, server, createOrUpdateUserArg{
		UserID:          userInfo.ID,
		UserLogin:       userInfo.Login,
		UserName:        userInfo.DisplayName,
		ProfileImageUrl: userInfo.ProfileImageURL,
		Token:           token,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// register eventsub
	appAccessToken, err := twClient.GetAppAccessToken(server.config.TwitchClientID, server.config.TwitchClientSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	eventsubPayload := twitch.EventsubSubscription{
		Type:      EventsubPayloadTypeStreamOnline,
		Version:   "1",
		Condition: twitch.EventsubSubscriptionCondition{BroadcasterUserID: userInfo.ID},
		Transport: twitch.EventsubSubscriptionTransport{
			Method:   "webhook",
			Callback: server.config.BeAddress + "/twitch/eventsub",
			Secret:   server.config.TwitchClientSecret,
		},
	}
	err = twClient.RegisterEventsub(appAccessToken.AccessToken, eventsubPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	eventsubPayload.Type = EventsubPayloadTypeChannelRedemption
	err = twClient.RegisterEventsub(appAccessToken.AccessToken, eventsubPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func createOrUpdateUser(ctx *gin.Context, server *Server, arg createOrUpdateUserArg) (loginUserResponse, error) {
	var output loginUserResponse

	_, err := server.store.GetUserByUserID(ctx, arg.UserID) // userId twitch
	if err != nil {
		if err == pgx.ErrNoRows {
			// create user
			argParams := db.CreateUserParams{
				UserID:          arg.UserID,
				UserLogin:       arg.UserLogin,
				UserName:        arg.UserName,
				ProfileImageUrl: arg.ProfileImageUrl,
			}

			user, err := server.store.CreateUser(ctx, argParams)
			if err != nil {
				// duplicate key error
				// should not be happened
				// because we already checked user_id in GetUserByUserID
				if strings.Contains(err.Error(), "duplicate key") {
					err = errors.New("unauthorized")
					return output, err
				}

				return output, err
			}

			rsp, err := createLoginUserResponse(ctx, server, user, arg.Token)
			if err != nil {
				return output, err
			}

			output = rsp
			return output, nil
		}

		return output, err
	}

	// update user
	updateUserArg := db.UpdateUserParams{
		UserID:          arg.UserID,
		UserLogin:       arg.UserLogin,
		UserName:        arg.UserName,
		ProfileImageUrl: arg.ProfileImageUrl,
	}

	user, err := server.store.UpdateUser(ctx, updateUserArg)
	if err != nil {
		return output, err
	}

	rsp, err := createLoginUserResponse(ctx, server, user, arg.Token)
	if err != nil {
		return output, err
	}

	output = rsp
	return output, nil
}

func createLoginUserResponse(ctx *gin.Context, server *Server, user db.User, token *twitch.OAuthToken) (loginUserResponse, error) {
	sessionID := uuid.New()

	// create token
	accessToken, _, err := server.tokenMaker.MakeToken(user.ID, sessionID, server.config.AccessTokenDuration)
	if err != nil {
		return loginUserResponse{}, err
	}

	// refresh token
	refreshToken, refreshTokenPayload, err := server.tokenMaker.MakeToken(user.ID, sessionID, server.config.RefreshTokenDuration*7)
	if err != nil {
		return loginUserResponse{}, err
	}

	// encrypt token
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return loginUserResponse{}, err
	}
	encryptedToken, err := util.Encrypt(string(tokenBytes), server.config.TokenSymmetricKey)
	if err != nil {
		return loginUserResponse{}, err
	}

	// create session
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:                   util.UUIDToUUID(sessionID),
		UserID:               user.ID,
		RefreshToken:         refreshToken,
		UserAgent:            ctx.Request.UserAgent(),
		ClientIp:             ctx.ClientIP(),
		IsBlocked:            false,
		ExpiresAt:            util.StringToTimestamp(refreshTokenPayload.ExpiredAt.Format(time.RFC3339)),
		EncryptedTwitchToken: encryptedToken,
	})
	if err != nil {
		return loginUserResponse{}, err
	}

	return loginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: userResponse{
			UserLogin:       user.UserLogin,
			UserName:        user.UserName,
			ProfileImageUrl: user.ProfileImageUrl,
		},
	}, nil
}

func (server *Server) refreshUser(ctx *gin.Context) {
	var req refreshUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		err = errors.New("invalid request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check refresh token
	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check session
	session, err := server.store.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if payload.UserID != session.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user")))
		return
	}
	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid refresh token")))
		return
	}
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("blocked")))
		return
	}

	accessToken, _, err := server.tokenMaker.MakeToken(payload.UserID, session.ID.Bytes, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

func (server *Server) createState(ctx *gin.Context) {
	state := util.RandomString(16)
	tempState[state] = true
	redirectURI := server.config.FeAddress + "/auth/login"

	url := "https://id.twitch.tv/oauth2/authorize?client_id=" + server.config.TwitchClientID + "&redirect_uri=" + redirectURI + "&response_type=code&scope=" + TwitchScope + "&state=" + state

	ctx.JSON(http.StatusOK, createStateResponse{
		URL: url,
	})
}

func (server *Server) logoutUser(ctx *gin.Context) {
	var req logoutUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check session
	session, err := server.store.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// delete session
	err = server.store.DeleteSession(ctx, session.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	delete(connectedTwitchClients, session.UserLogin)
	delete(connectedWsClients, session.UserLogin)

	// delete twitch chat
	twitch.Disconnect(session.UserLogin)

	ctx.JSON(http.StatusOK, nil)
}
