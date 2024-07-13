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
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserConfig(ctx context.Context, id int64) (UserConfig, error) {
	row := q.db.QueryRow(ctx, getUserConfig, id)
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

const listUserConfigs = `-- name: ListUserConfigs :many
SELECT id, user_id, config_type, value, created_at, updated_at FROM user_configs
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUserConfigsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUserConfigs(ctx context.Context, arg ListUserConfigsParams) ([]UserConfig, error) {
	rows, err := q.db.Query(ctx, listUserConfigs, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserConfig
	for rows.Next() {
		var i UserConfig
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ConfigType,
			&i.Value,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
