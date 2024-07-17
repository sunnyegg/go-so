package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

type createUserRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	UserLogin       string `json:"user_login" binding:"required"`
	UserName        string `json:"user_name" binding:"required"`
	ProfileImageUrl string `json:"profile_image_url"`
	Token           string `json:"token"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash token
	hashedToken, err := util.HashToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// TODO: create jwt token

	// check user login
	// if not exists, createUser
	// else updateUser
	user, err := server.store.GetUserByUserID(ctx, req.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// create user
			arg := db.CreateUserParams{
				UserLogin:       user.UserLogin,
				UserName:        user.UserName,
				ProfileImageUrl: user.ProfileImageUrl,
				Token:           hashedToken,
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

			ctx.JSON(http.StatusOK, user)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update user
	arg := db.UpdateUserParams{
		UserLogin:       req.UserLogin,
		UserName:        req.UserName,
		ProfileImageUrl: req.ProfileImageUrl,
		Token:           hashedToken,
	}

	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
