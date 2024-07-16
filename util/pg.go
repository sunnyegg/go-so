package util

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// convert string to pgtype.Text
func StringToText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

// convert string to pgtype.Timestamp
func StringToTimestamp(s string) pgtype.Timestamptz {
	// convert string to time.Time
	// time rfc3339 example: 2006-01-02T15:04:05Z07:00
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Printf("failed to parse time: %s", err)
	}
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
