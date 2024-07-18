// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: sessions.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions
(
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
)
VALUES
($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
`

type CreateSessionParams struct {
	ID           pgtype.UUID        `json:"id"`
	UserID       int64              `json:"user_id"`
	RefreshToken string             `json:"refresh_token"`
	UserAgent    string             `json:"user_agent"`
	ClientIp     string             `json:"client_ip"`
	IsBlocked    bool               `json:"is_blocked"`
	ExpiresAt    pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions
WHERE id = $1 AND user_id = $2
LIMIT 1
`

type GetSessionParams struct {
	ID     pgtype.UUID `json:"id"`
	UserID int64       `json:"user_id"`
}

func (q *Queries) GetSession(ctx context.Context, arg GetSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, arg.ID, arg.UserID)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
