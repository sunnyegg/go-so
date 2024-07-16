package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/sunnyegg/go-so/sqlc"
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

	arg := db.CreateAttendanceMemberParams{
		StreamID:  req.StreamID,
		Username:  req.Username,
		IsShouted: false,
		PresentAt: util.StringToTimestamp(req.PresentAt),
	}

	attendanceMember, err := server.queries.CreateAttendanceMember(ctx, arg)
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
