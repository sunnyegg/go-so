package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	plainText := "Hello, World!"
	cipherText, err := encrypt(plainText, "examplekey123456")
	require.NoError(t, err)
	require.NotEmpty(t, cipherText)

	decryptedText, err := decrypt(cipherText, "examplekey123456")
	require.NoError(t, err)
	require.Equal(t, plainText, decryptedText)
}
