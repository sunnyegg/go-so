package cron

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

func ValidateToken(ctx context.Context, store db.Store, config util.Config) func() {
	twClient := twitch.NewClient(config.TwitchClientID, config.TwitchClientSecret, config.FeAddress)

	return func() {
		// get sessions
		sessions, err := store.ListSession(ctx)
		if err != nil {
			return
		}

		// loop sessions and validate token
		for _, session := range sessions {
			sessionID := uuid.UUID(session.ID.Bytes)
			log.Printf("validating session %s", sessionID)

			decryptedToken, err := util.Decrypt(session.EncryptedTwitchToken, config.TokenSymmetricKey)
			if err != nil {
				log.Println("failed to decrypt token", err)
				return
			}
			// unmarshal decrypted token
			// string to []byte
			token := []byte(decryptedToken)
			payload := twitch.OAuthToken{}
			err = json.Unmarshal(token, &payload)
			if err != nil {
				log.Println("failed to unmarshal token", err)
				return
			}

			// validate token
			_, err = twClient.ValidateOAuthToken(payload.AccessToken)
			if err != nil {
				log.Println("failed to validate token", err)
				return
			}

			// token's session is valid
			log.Printf("session %s token is valid", sessionID)
		}
	}
}
