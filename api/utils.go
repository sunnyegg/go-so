package api

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/token"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

func decryptHeader(ctx *gin.Context, server *Server) (*twitch.OAuthToken, *token.Payload, error) {
	// get sessionid
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// session
	session, err := server.store.GetSession(ctx, db.GetSessionParams{
		ID:     util.UUIDToUUID(authPayload.SessionID),
		UserID: authPayload.UserID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, authPayload, errors.New("unauthorized")
		}

		return nil, authPayload, err
	}

	// decrypt token
	tokenBytes, err := util.Decrypt(session.EncryptedTwitchToken, server.config.TokenSymmetricKey)
	if err != nil {
		return nil, authPayload, err
	}
	token := []byte(tokenBytes)
	payload := twitch.OAuthToken{}
	err = json.Unmarshal(token, &payload)
	if err != nil {
		return nil, authPayload, err
	}

	return &payload, authPayload, nil
}
