package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/sunnyegg/go-so/util"
)

func createRandomSession(t *testing.T, userID *int64) Session {
	var user User
	if userID == nil {
		user = createRandomUser(t)
	} else {
		user.ID = *userID
	}

	arg := CreateSessionParams{
		ID:                   util.UUIDToUUID(uuid.New()),
		UserID:               user.ID,
		RefreshToken:         util.RandomString(32),
		UserAgent:            util.RandomString(10),
		ClientIp:             util.RandomString(10),
		IsBlocked:            false,
		ExpiresAt:            util.StringToTimestamp(time.Now().Add(time.Minute).Format(time.RFC3339)),
		EncryptedTwitchToken: util.RandomString(10),
	}

	session, err := testStore.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.NotEmpty(t, session.ID)
	require.Equal(t, arg.UserID, session.UserID)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.Equal(t, arg.EncryptedTwitchToken, session.EncryptedTwitchToken)

	require.WithinDuration(t, arg.ExpiresAt.Time, session.ExpiresAt.Time, time.Second)

	return session
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t, nil)
}

func TestGetSession(t *testing.T) {
	session1 := createRandomSession(t, nil)
	arg := GetSessionParams{
		ID:     session1.ID,
		UserID: session1.UserID,
	}
	session2, err := testStore.GetSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.UserID, session2.UserID)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.Equal(t, session1.EncryptedTwitchToken, session2.EncryptedTwitchToken)
}

func TestListSession(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomSession(t, nil)
	}

	sessions, err := testStore.ListSession(context.Background())
	require.NoError(t, err)
	require.Greater(t, len(sessions), 9)

	for _, session := range sessions {
		require.NotEmpty(t, session)
	}
}

func TestUpdateSession(t *testing.T) {
	session1 := createRandomSession(t, nil)
	arg := UpdateSessionParams{
		ID:                   session1.ID,
		EncryptedTwitchToken: util.RandomString(10),
	}
	err := testStore.UpdateSession(context.Background(), arg)
	require.NoError(t, err)

	session2, err := testStore.GetSession(context.Background(), GetSessionParams{
		ID:     session1.ID,
		UserID: session1.UserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, arg.EncryptedTwitchToken, session2.EncryptedTwitchToken)
}
