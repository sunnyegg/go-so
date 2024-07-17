package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/sunnyegg/go-so/util"
)

func createRandomAttendanceMember(t *testing.T, streamID *int64) AttendanceMember {
	var stream Stream
	if streamID == nil {
		stream = createRandomStream(t, nil)
	} else {
		stream.ID = *streamID
	}

	arg := CreateAttendanceMemberParams{
		StreamID:  stream.ID,
		Username:  util.RandomString(10),
		IsShouted: false,
		PresentAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	attendanceMember, err := testStore.CreateAttendanceMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, attendanceMember)

	require.Equal(t, arg.StreamID, attendanceMember.StreamID)
	require.Equal(t, arg.Username, attendanceMember.Username)
	require.Equal(t, arg.IsShouted, attendanceMember.IsShouted)

	require.NotZero(t, attendanceMember.PresentAt)

	return attendanceMember
}

func TestCreateAttendanceMember(t *testing.T) {
	createRandomAttendanceMember(t, nil)
}
