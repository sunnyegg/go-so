// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user_configs.sql

package db

import (
	"context"
)

const createUserConfig = `-- name: CreateUserConfig :one
INSERT INTO user_configs
(
  user_id,
  config_type,
  "value"
)
VALUES
($1, $2, $3)
RETURNING id, user_id, config_type, value, created_at, updated_at
`

type CreateUserConfigParams struct {
	UserID     int64       `json:"user_id"`
	ConfigType ConfigTypes `json:"config_type"`
	Value      string      `json:"value"`
}

func (q *Queries) CreateUserConfig(ctx context.Context, arg CreateUserConfigParams) (UserConfig, error) {
	row := q.db.QueryRow(ctx, createUserConfig, arg.UserID, arg.ConfigType, arg.Value)
	var i UserConfig
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ConfigType,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserConfig = `-- name: DeleteUserConfig :exec
DELETE FROM user_configs
WHERE id = $1
`

func (q *Queries) DeleteUserConfig(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUserConfig, id)
	return err
}

const getUserConfig = `-- name: GetUserConfig :one
SELECT id, user_id, config_type, value, created_at, updated_at FROM user_configs
WHERE user_id = $1 AND config_type = $2
LIMIT 1
`

type GetUserConfigParams struct {
	UserID     int64       `json:"user_id"`
	ConfigType ConfigTypes `json:"config_type"`
}

func (q *Queries) GetUserConfig(ctx context.Context, arg GetUserConfigParams) (UserConfig, error) {
	row := q.db.QueryRow(ctx, getUserConfig, arg.UserID, arg.ConfigType)
	var i UserConfig
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ConfigType,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserConfig = `-- name: UpdateUserConfig :one
UPDATE user_configs
SET "value" = $2, updated_at = now()
WHERE id = $1
RETURNING id, user_id, config_type, value, created_at, updated_at
`

type UpdateUserConfigParams struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

func (q *Queries) UpdateUserConfig(ctx context.Context, arg UpdateUserConfigParams) (UserConfig, error) {
	row := q.db.QueryRow(ctx, updateUserConfig, arg.ID, arg.Value)
	var i UserConfig
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ConfigType,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
