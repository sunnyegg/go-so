package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAttendance(t *testing.T) Attendance {
	stream := createRandomStream(t)

	arg := CreateAttendanceParams{
		UserID:   stream.UserID,
		StreamID: stream.ID,
	}

	attendance, err := testQueries.CreateAttendance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, attendance)

	require.Equal(t, arg.UserID, attendance.UserID)
	require.Equal(t, arg.StreamID, attendance.StreamID)

	require.NotZero(t, attendance.ID)
	require.NotZero(t, attendance.CreatedAt)

	return attendance
}

func TestCreateAttendance(t *testing.T) {
	createRandomAttendance(t)
}

func TestGetAttendance(t *testing.T) {
	attendance1 := createRandomAttendance(t)
	attendance2, err := testQueries.GetAttendance(context.Background(), attendance1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, attendance2)

	require.Equal(t, attendance1.UserID, attendance2.UserID)
	require.Equal(t, attendance1.StreamID, attendance2.StreamID)
	require.WithinDuration(t, attendance1.CreatedAt.Time, attendance2.CreatedAt.Time, time.Second)
}

func TestDeleteAttendance(t *testing.T) {
	attendance1 := createRandomAttendance(t)
	err := testQueries.DeleteAttendance(context.Background(), attendance1.ID)
	require.NoError(t, err)

	attendance2, err := testQueries.GetAttendance(context.Background(), attendance1.ID)
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
	require.Empty(t, attendance2)
}

func TestListAttendances(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAttendance(t)
	}

	arg := ListAttendancesParams{
		Limit:  5,
		Offset: 5,
	}

	attendances, err := testQueries.ListAttendances(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, attendances, 5)

	for _, attendance := range attendances {
		require.NotEmpty(t, attendance)
	}
}
