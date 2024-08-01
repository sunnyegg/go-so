package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

type createStreamResponse struct {
	ID        int64  `json:"id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	GameName  string `json:"game_name" binding:"required"`
	StartedAt string `json:"started_at" binding:"required"`
}

func (server *Server) createStream(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// get session
	session, err := server.store.GetSession(ctx, db.GetSessionParams{
		ID:     util.UUIDToUUID(authPayload.SessionID),
		UserID: authPayload.UserID,
	})
	if err != nil {
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

	// get stream info
	twClient := twitch.NewClient(server.config.TwitchClientID, server.config.TwitchClientSecret, server.config.RedirectURI)
	streamInfo, err := twClient.GetStreamInfo(payload.AccessToken, session.UserID_2)
	if err != nil {
		if err.Error() == "stream not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateStreamParams{
		UserID:    authPayload.UserID,
		Title:     streamInfo.Title,
		GameName:  streamInfo.GameName,
		StartedAt: util.StringToTimestamp(streamInfo.StartedAt),
		CreatedBy: authPayload.UserID,
	}

	stream, err := server.store.CreateStream(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createStreamResponse{
		ID:        stream.ID,
		Title:     stream.Title,
		GameName:  stream.GameName,
		StartedAt: stream.StartedAt.Time.String(),
	})
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetStreamParams{
		ID:     req.ID,
		UserID: authPayload.UserID,
	}

	stream, err := server.store.GetStream(ctx, arg)
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
}

func (server *Server) listStream(ctx *gin.Context) {
	var req listStreamRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListStreamsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		UserID: authPayload.UserID,
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

func (server *Server) getStreamAttendanceMember(ctx *gin.Context) {
	var req getStreamAttendanceMembersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.GetStreamAttendanceMembersParams{
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
		StreamID: req.StreamID,
		UserID:   authPayload.UserID,
	}

	attendanceMembers, err := server.store.GetStreamAttendanceMembers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, attendanceMembers)
}
