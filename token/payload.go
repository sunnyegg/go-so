package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userID int64, sessionID uuid.UUID, duration time.Duration) (*Payload, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        uid,
		UserID:    userID,
		SessionID: sessionID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if payload.ID == uuid.Nil {
		return ErrInvalidToken
	}

	if payload.IssuedAt.IsZero() {
		return ErrInvalidToken
	}

	if payload.ExpiredAt.IsZero() {
		return ErrInvalidToken
	}

	if payload.ExpiredAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}
