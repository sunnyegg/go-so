package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
)

type createUserConfigRequest struct {
	ConfigType db.ConfigTypes `json:"config_type" binding:"required"`
	Value      string         `json:"value" binding:"required"`
}

func (server *Server) createUserConfig(ctx *gin.Context) {
	var req createUserConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	userConfig, err := server.store.CreateUserConfig(ctx, db.CreateUserConfigParams{
		UserID:     authPayload.UserID,
		ConfigType: req.ConfigType,
		Value:      req.Value,
	})
	if err != nil {
		if strings.Contains(err.Error(), "invalid input value for enum config_types") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid config type")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userConfig)
}

type getUserConfigRequest struct {
	ConfigType db.ConfigTypes `form:"config_type" binding:"required,min=1"`
}

func (server *Server) getUserConfig(ctx *gin.Context) {
	var req getUserConfigRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	userConfig, err := server.store.GetUserConfig(ctx, db.GetUserConfigParams{
		UserID:     authPayload.UserID,
		ConfigType: req.ConfigType,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		if strings.Contains(err.Error(), "invalid input value for enum config_types") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid config type")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userConfig)
}
