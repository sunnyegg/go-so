package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashToken(t *testing.T) {
	token1, err := HashToken("test")
	require.NoError(t, err)
	require.NotEmpty(t, token1)

	err = VerifyToken("test", token1)
	require.NoError(t, err)

	err = VerifyToken("wrong", token1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// if using the same string, it should return different hash
	token2, err := HashToken("test")
	require.NoError(t, err)
	require.NotEqual(t, token1, token2)
}
