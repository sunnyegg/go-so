package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

type loginUserRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	UserLogin       string `json:"user_login" binding:"required"`
	UserName        string `json:"user_name" binding:"required"`
	ProfileImageUrl string `json:"profile_image_url"`
	Token           string `json:"token" binding:"required"`
}

type userResponse struct {
	UserLogin       string `json:"user_login"`
	UserName        string `json:"user_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type loginUserResponse struct {
	SessionID    uuid.UUID    `json:"session_id"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check user login
	// if not exists, createUser
	// else updateUser
	_, err := server.store.GetUserByUserID(ctx, req.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// create user
			arg := db.CreateUserParams{
				UserID:          req.UserID,
				UserLogin:       req.UserLogin,
				UserName:        req.UserName,
				ProfileImageUrl: req.ProfileImageUrl,
			}

			user, err := server.store.CreateUser(ctx, arg)
			if err != nil {
				// duplicate key error
				// should not be happened
				// because we already checked user_id in GetUserByUserID
				if strings.Contains(err.Error(), "duplicate key") {
					ctx.JSON(http.StatusForbidden, errorResponse(errors.New("unauthorized")))
					return
				}

				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			rsp, err := createLoginUserResponse(ctx, user, server)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusOK, rsp)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update user
	arg := db.UpdateUserParams{
		UserID:          req.UserID,
		UserLogin:       req.UserLogin,
		UserName:        req.UserName,
		ProfileImageUrl: req.ProfileImageUrl,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := createLoginUserResponse(ctx, user, server)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func createLoginUserResponse(ctx *gin.Context, user db.User, server *Server) (loginUserResponse, error) {
	// create token
	accessToken, _, err := server.tokenMaker.MakeToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		return loginUserResponse{}, err
	}

	// refresh token
	refreshToken, payload, err := server.tokenMaker.MakeToken(user.ID, server.config.RefreshTokenDuration*7)
	if err != nil {
		return loginUserResponse{}, err
	}

	// create session
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           util.UUIDToUUID(payload.ID),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    util.StringToTimestamp(payload.ExpiredAt.Format(time.RFC3339)),
	})
	if err != nil {
		return loginUserResponse{}, err
	}

	return loginUserResponse{
		SessionID:    payload.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: userResponse{
			UserLogin:       user.UserLogin,
			UserName:        user.UserName,
			ProfileImageUrl: user.ProfileImageUrl,
		},
	}, nil
}

type refreshUserRequest struct {
	SessionID    uuid.UUID `json:"session_id" binding:"required"`
	RefreshToken string    `json:"refresh_token" binding:"required"`
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
	if payload.ID != req.SessionID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
		return
	}
	if payload.ExpiredAt.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("expired")))
		return
	}

	// check session
	session, err := server.store.GetSession(ctx, db.GetSessionParams{
		ID:     util.UUIDToUUID(payload.ID),
		UserID: payload.UserID,
	})
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

	accessToken, _, err := server.tokenMaker.MakeToken(payload.UserID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}
