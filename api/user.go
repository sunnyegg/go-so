package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
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

	// TODO: hash token

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
				Token:           user.Token,
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
		Token:           req.Token,
	}

	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}
