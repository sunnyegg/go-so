package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/sunnyegg/go-so/util"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		UserID:          util.RandomUserID(),
		UserLogin:       util.RandomString(8),
		UserName:        util.RandomString(8),
		ProfileImageUrl: util.RandomString(16),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserID, user.UserID)
	require.Equal(t, arg.UserLogin, user.UserLogin)
	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.ProfileImageUrl, user.ProfileImageUrl)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserLogin, user2.UserLogin)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.ProfileImageUrl, user2.ProfileImageUrl)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserParams{
		UserID:          user1.UserID,
		UserLogin:       user1.UserLogin,
		UserName:        util.RandomString(8),
		ProfileImageUrl: util.RandomString(16),
	}
	user2, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.UserID, user2.UserID)
	require.Equal(t, arg.UserLogin, user2.UserLogin)
	require.Equal(t, arg.UserName, user2.UserName)
	require.Equal(t, arg.ProfileImageUrl, user2.ProfileImageUrl)
	require.NotZero(t, user2.UpdatedAt)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testStore.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testStore.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
