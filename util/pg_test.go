package util

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStringToText(t *testing.T) {
	s := RandomString(10)
	text := StringToText(s)
	require.NotEmpty(t, text)
	require.Equal(t, s, text.String)
}

func TestStringToTimestamp(t *testing.T) {
	s := time.Now().Format(time.RFC3339)
	timestamp := StringToTimestamp(s)
	require.NotEmpty(t, timestamp)
	require.Equal(t, s, timestamp.Time.Format(time.RFC3339))
}

func TestUUIDToUUID(t *testing.T) {
	s := uuid.New().String()
	pgUuid := UUIDToUUID(uuid.MustParse(s))
	require.NotEmpty(t, pgUuid)
}
