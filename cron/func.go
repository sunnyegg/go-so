package cron

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/sunnyegg/go-so/channel"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/twitch"
	"github.com/sunnyegg/go-so/util"
)

func ValidateToken(ctx context.Context, store db.Store, config util.Config) func() {
	twClient := twitch.NewClient(config.TwitchClientID, config.TwitchClientSecret, config.FeAddress)
	ch := channel.NewChannel(channel.ChannelGeneral)

	return func() {
		// get sessions
		sessions, err := store.ListSession(ctx)
		if err != nil {
			return
		}
		log.Printf("refreshing %d sessions", len(sessions))

		// loop sessions and validate token
		for i, session := range sessions {
			sessionID := uuid.UUID(session.ID.Bytes)
			log.Printf("validating session[%d] %s", i, sessionID)

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
			err = twClient.ValidateOAuthToken(payload.AccessToken)
			if err != nil {
				if err.Error() != twitch.ErrExpiredToken {
					log.Println("failed to validate token", err)
					return
				}

				log.Println("refreshing token...")
				refreshedToken, err := twClient.RefreshOAuthToken(payload.RefreshToken)
				if err != nil {
					log.Println("failed to refresh token", err)
					return
				}

				// send token to all connected chat clients
				go ch.Send(map[string]string{
					"channel": session.UserLogin,
					"token":   refreshedToken.AccessToken,
				})

				// encrypt token
				log.Println("encrypting token...")
				tokenBytes, err := json.Marshal(refreshedToken)
				if err != nil {
					log.Println("failed to marshal token", err)
					return
				}
				encryptedToken, err := util.Encrypt(string(tokenBytes), config.TokenSymmetricKey)
				if err != nil {
					log.Println("failed to encrypt token", err)
					return
				}

				// update session
				log.Println("updating session...")
				err = store.UpdateSession(ctx, db.UpdateSessionParams{
					ID:                   session.ID,
					EncryptedTwitchToken: encryptedToken,
				})
				if err != nil {
					log.Println("failed to update session", err)
					return
				}

				log.Printf("session[%d] %s token is refreshed", i, sessionID)
			}

			// token's session is valid
			log.Printf("session[%d] %s token is valid", i, sessionID)
		}
	}
}
