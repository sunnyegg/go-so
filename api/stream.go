package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

type createStreamRequest struct {
	UserID    int64  `json:"user_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	GameName  string `json:"game_name" binding:"required"`
	StartedAt string `json:"started_at" binding:"required"`
}

func (server *Server) createStream(ctx *gin.Context) {
	var req createStreamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateStreamParams{
		UserID:    req.UserID,
		Title:     req.Title,
		GameName:  req.GameName,
		StartedAt: util.StringToTimestamp(req.StartedAt),
		CreatedBy: req.UserID, // TODO: get from token
	}

	stream, err := server.store.CreateStream(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, stream)
}

type getStreamRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getStream(ctx *gin.Context) {
	var req getStreamRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	stream, err := server.store.GetStream(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, stream)
}

type listStreamRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	UserID   int64 `form:"user_id" binding:"required,min=1"`
}

// TODO: filter by user_id in token
func (server *Server) listStream(ctx *gin.Context) {
	var req listStreamRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListStreamsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		UserID: req.UserID,
	}

	streams, err := server.store.ListStreams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, streams)
}

type getStreamAttendanceMembersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	StreamID int64 `form:"stream_id" binding:"required,min=1"`
}

// TODO: filter by user_id in token
func (server *Server) getStreamAttendanceMember(ctx *gin.Context) {
	var req getStreamAttendanceMembersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetStreamAttendanceMembersParams{
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
		StreamID: req.StreamID,
	}

	attendanceMembers, err := server.store.GetStreamAttendanceMembers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, attendanceMembers)
}
