// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users
(
  user_id,
  user_login,
  user_name,
  profile_image_url
)
VALUES
($1, $2, $3, $4)
RETURNING id, user_id, user_login, user_name, profile_image_url, created_at, updated_at
`

type CreateUserParams struct {
	UserID          string      `json:"user_id"`
	UserLogin       string      `json:"user_login"`
	UserName        string      `json:"user_name"`
	ProfileImageUrl pgtype.Text `json:"profile_image_url"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UserID,
		arg.UserLogin,
		arg.UserName,
		arg.ProfileImageUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserLogin,
		&i.UserName,
		&i.ProfileImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, user_id, user_login, user_name, profile_image_url, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserLogin,
		&i.UserName,
		&i.ProfileImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, user_id, user_login, user_name, profile_image_url, created_at, updated_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.UserLogin,
			&i.UserName,
			&i.ProfileImageUrl,
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

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET user_login = $2, user_name = $3, profile_image_url = $4, updated_at = now()
WHERE id = $1
RETURNING id, user_id, user_login, user_name, profile_image_url, created_at, updated_at
`

type UpdateUserParams struct {
	ID              int64       `json:"id"`
	UserLogin       string      `json:"user_login"`
	UserName        string      `json:"user_name"`
	ProfileImageUrl pgtype.Text `json:"profile_image_url"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.UserLogin,
		arg.UserName,
		arg.ProfileImageUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserLogin,
		&i.UserName,
		&i.ProfileImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
