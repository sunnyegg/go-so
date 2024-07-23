package util

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// convert string to pgtype.Text
func StringToText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

// convert string timestamp to pgtype.Timestamp
func StringToTimestamp(s string) pgtype.Timestamptz {
	// convert string to time.Time
	// time rfc3339 example: 2006-01-02T15:04:05Z07:00
	t, _ := time.Parse(time.RFC3339, s)
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// convert uuid.UUID to pgtype.UUID
func UUIDToUUID(s uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: s,
		Valid: true,
	}
}
