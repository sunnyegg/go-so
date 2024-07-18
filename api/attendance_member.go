package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/util"
)

type createAttendanceMemberRequest struct {
	StreamID  int64  `json:"stream_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	PresentAt string `json:"present_at" binding:"required"`
}

func (server *Server) createAttendanceMember(ctx *gin.Context) {
	var req createAttendanceMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// check is stream related to user
	_, err := server.store.GetStream(ctx, db.GetStreamParams{
		ID:     req.StreamID,
		UserID: authPayload.UserID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateAttendanceMemberParams{
		StreamID:  req.StreamID,
		Username:  req.Username,
		IsShouted: false,
		PresentAt: util.StringToTimestamp(req.PresentAt),
	}

	attendanceMember, err := server.store.CreateAttendanceMember(ctx, arg)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("member exists")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, attendanceMember)
}
