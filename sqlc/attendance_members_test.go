package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/sunnyegg/go-so/util"
)

func createRandomAttendanceMember(t *testing.T, attendanceID *int64) AttendanceMember {
	var attendance Attendance
	if attendanceID == nil {
		attendance = createRandomAttendance(t)
	} else {
		attendance.ID = *attendanceID
	}

	arg := CreateAttendanceMemberParams{
		AttendanceID: attendance.ID,
		Username:     util.RandomString(10),
		IsShouted:    false,
		PresentAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	attendanceMember, err := testQueries.CreateAttendanceMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, attendanceMember)

	require.Equal(t, arg.AttendanceID, attendanceMember.AttendanceID)
	require.Equal(t, arg.Username, attendanceMember.Username)
	require.Equal(t, arg.IsShouted, attendanceMember.IsShouted)

	require.NotZero(t, attendanceMember.ID)
	require.NotZero(t, attendanceMember.CreatedAt)
	require.NotZero(t, attendanceMember.PresentAt)

	return attendanceMember
}

func TestCreateAttendanceMember(t *testing.T) {
	createRandomAttendanceMember(t, nil)
}

func TestGetAttendanceMember(t *testing.T) {
	attendanceMember1 := createRandomAttendanceMember(t, nil)
	attendanceMember2, err := testQueries.GetAttendanceMember(context.Background(), attendanceMember1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, attendanceMember2)

	require.Equal(t, attendanceMember1.AttendanceID, attendanceMember2.AttendanceID)
	require.Equal(t, attendanceMember1.Username, attendanceMember2.Username)
	require.Equal(t, attendanceMember1.IsShouted, attendanceMember2.IsShouted)
	require.WithinDuration(t, attendanceMember1.PresentAt.Time, attendanceMember2.PresentAt.Time, time.Second)
	require.WithinDuration(t, attendanceMember1.CreatedAt.Time, attendanceMember2.CreatedAt.Time, time.Second)
}

func TestListAttendanceMembers(t *testing.T) {
	attendance := createRandomAttendanceMember(t, nil)
	for i := 0; i < 9; i++ {
		createRandomAttendanceMember(t, &attendance.ID)
	}

	arg := ListAttendanceMembersParams{
		Limit:  5,
		Offset: 5,
	}

	attendanceMembers, err := testQueries.ListAttendanceMembers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, attendanceMembers, 5)

	for _, attendanceMember := range attendanceMembers {
		require.NotEmpty(t, attendanceMember)
	}
}
