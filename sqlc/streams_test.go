package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/sunnyegg/go-so/util"
)

func createRandomStream(t *testing.T) Stream {
	user := createRandomUser(t)

	arg := CreateStreamParams{
		UserID:   user.ID,
		Title:    util.RandomString(10),
		GameName: util.RandomString(20),
		StartedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedBy: user.ID,
	}

	stream, err := testQueries.CreateStream(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, stream)

	require.Equal(t, arg.UserID, stream.UserID)
	require.Equal(t, arg.Title, stream.Title)
	require.Equal(t, arg.GameName, stream.GameName)
	require.Equal(t, arg.CreatedBy, stream.CreatedBy)

	require.NotZero(t, stream.ID)
	require.NotZero(t, arg.StartedAt, stream.StartedAt)
	require.NotZero(t, stream.CreatedAt)

	return stream
}

func TestCreateStream(t *testing.T) {
	createRandomStream(t)
}

func TestGetStream(t *testing.T) {
	stream1 := createRandomStream(t)
	stream2, err := testQueries.GetStream(context.Background(), stream1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, stream2)

	require.Equal(t, stream1.UserID, stream2.UserID)
	require.Equal(t, stream1.Title, stream2.Title)
	require.Equal(t, stream1.GameName, stream2.GameName)
	require.Equal(t, stream1.CreatedBy, stream2.CreatedBy)
	require.WithinDuration(t, stream1.StartedAt.Time, stream2.StartedAt.Time, time.Second)
	require.WithinDuration(t, stream1.CreatedAt.Time, stream2.CreatedAt.Time, time.Second)
}

func TestDeleteStream(t *testing.T) {
	stream1 := createRandomStream(t)
	err := testQueries.DeleteStream(context.Background(), stream1.ID)
	require.NoError(t, err)

	stream2, err := testQueries.GetStream(context.Background(), stream1.ID)
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
	require.Empty(t, stream2)
}

func TestListStreams(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomStream(t)
	}

	arg := ListStreamsParams{
		Limit:  5,
		Offset: 5,
	}

	streams, err := testQueries.ListStreams(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, streams, 5)

	for _, stream := range streams {
		require.NotEmpty(t, stream)
	}
}
