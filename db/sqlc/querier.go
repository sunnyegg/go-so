// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CreateAttendanceMember(ctx context.Context, arg CreateAttendanceMemberParams) (AttendanceMember, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateStream(ctx context.Context, arg CreateStreamParams) (Stream, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserConfig(ctx context.Context, arg CreateUserConfigParams) (UserConfig, error)
	DeleteStream(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteUserConfig(ctx context.Context, id int64) error
	GetSession(ctx context.Context, arg GetSessionParams) (Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	GetSessionByUserID(ctx context.Context, userID string) (GetSessionByUserIDRow, error)
	GetStream(ctx context.Context, arg GetStreamParams) (GetStreamRow, error)
	GetStreamAttendanceMembers(ctx context.Context, arg GetStreamAttendanceMembersParams) ([]GetStreamAttendanceMembersRow, error)
	GetUser(ctx context.Context, id int64) (GetUserRow, error)
	GetUserByUserID(ctx context.Context, userID string) (User, error)
	GetUserConfig(ctx context.Context, arg GetUserConfigParams) (UserConfig, error)
	ListSession(ctx context.Context) ([]Session, error)
	ListStreams(ctx context.Context, arg ListStreamsParams) ([]ListStreamsRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateSession(ctx context.Context, arg UpdateSessionParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserConfig(ctx context.Context, arg UpdateUserConfigParams) (UserConfig, error)
}

var _ Querier = (*Queries)(nil)
