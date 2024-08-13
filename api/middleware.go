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

const (
	authorizationHeaderKey  = "Authorization"
	authorizationPrefixKey  = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(authorizationHeaderKey)
		if len(token) == 0 {
			err := errors.New("authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		field := strings.Fields(token)
		if len(field) < 2 {
			err := errors.New("invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(field[0])
		if authorizationType != authorizationPrefixKey {
			err := errors.New("invalid authorization type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := field[1]
		payload, err := server.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// check if user is blocked
		session, err := server.store.GetSession(c, db.GetSessionParams{
			ID:     util.UUIDToUUID(payload.SessionID),
			UserID: payload.UserID,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				err := errors.New("session not found")
				c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if session.IsBlocked {
			err := errors.New("session is blocked")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
