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

func createRandomStream(t *testing.T, userID *int64) Stream {
	var user User
	if userID == nil {
		user = createRandomUser(t)
	} else {
		user.ID = *userID
	}

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

	stream, err := testStore.CreateStream(context.Background(), arg)
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
	createRandomStream(t, nil)
}

func TestGetStream(t *testing.T) {
	stream1 := createRandomStream(t, nil)
	arg := GetStreamParams{
		ID:     stream1.ID,
		UserID: stream1.UserID,
	}
	stream2, err := testStore.GetStream(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, stream2)

	require.Equal(t, stream1.Title, stream2.Title)
	require.Equal(t, stream1.GameName, stream2.GameName)
	require.NotZero(t, stream2.CreatedBy)
	require.WithinDuration(t, stream1.StartedAt.Time, stream2.StartedAt.Time, time.Second)
}

func TestDeleteStream(t *testing.T) {
	stream1 := createRandomStream(t, nil)
	err := testStore.DeleteStream(context.Background(), stream1.ID)
	require.NoError(t, err)

	arg := GetStreamParams{
		ID:     stream1.ID,
		UserID: stream1.UserID,
	}
	stream2, err := testStore.GetStream(context.Background(), arg)
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
	require.Empty(t, stream2)
}

func TestListStreams(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomStream(t, &user.ID)
	}

	arg := ListStreamsParams{
		Limit:  5,
		Offset: 5,
		UserID: user.ID,
	}

	streams, err := testStore.ListStreams(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, streams, 5)

	for _, stream := range streams {
		require.NotEmpty(t, stream)
	}
}

func TestGetStreamAttendanceMembers(t *testing.T) {
	stream1 := createRandomStream(t, nil)
	createRandomAttendanceMember(t, &stream1.ID)
	createRandomAttendanceMember(t, &stream1.ID)

	arg := GetStreamAttendanceMembersParams{
		Limit:    5,
		Offset:   0,
		StreamID: stream1.ID,
		UserID:   stream1.UserID,
	}

	attendanceMembers, err := testStore.GetStreamAttendanceMembers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, attendanceMembers, 2)

	for _, attendanceMember := range attendanceMembers {
		require.NotEmpty(t, attendanceMember)
	}
}
