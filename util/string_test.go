package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	s := RandomUserID()
	i, err := ParseStringToInt64(s)
	require.NoError(t, err)
	require.NotEmpty(t, i)

	i32, err := ParseStringToInt32(s)
	require.NoError(t, err)
	require.NotEmpty(t, i32)
}

func TestEmptyParseString(t *testing.T) {
	s := ""
	i, err := ParseStringToInt64(s)
	require.Error(t, err, errors.New("empty string"))
	require.Equal(t, int64(0), i)

	i32, err := ParseStringToInt32(s)
	require.Error(t, err, errors.New("empty string"))
	require.Equal(t, int32(0), i32)
}

func TestInvalidParseString(t *testing.T) {
	s := RandomString(10)
	i, err := ParseStringToInt64(s)
	require.Error(t, err)
	require.Equal(t, int64(0), i)

	i32, err := ParseStringToInt32(s)
	require.Error(t, err)
	require.Equal(t, int32(0), i32)
}
