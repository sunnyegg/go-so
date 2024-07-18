package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
)

type loginUserRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	UserLogin       string `json:"user_login" binding:"required"`
	UserName        string `json:"user_name" binding:"required"`
	ProfileImageUrl string `json:"profile_image_url"`
	Token           string `json:"token"`
}

type userResponse struct {
	UserLogin       string `json:"user_login"`
	UserName        string `json:"user_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
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

			rsp, err := createLoginUserResponse(user, server)
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

	rsp, err := createLoginUserResponse(user, server)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func createLoginUserResponse(user db.User, server *Server) (loginUserResponse, error) {
	// create token
	accessToken, err := server.tokenMaker.MakeToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		return loginUserResponse{}, err
	}

	return loginUserResponse{
		AccessToken: accessToken,
		User: userResponse{
			UserLogin:       user.UserLogin,
			UserName:        user.UserName,
			ProfileImageUrl: user.ProfileImageUrl,
		},
	}, nil
}
