package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	plainText := "Hello, World!"
	cipherText, err := Encrypt(plainText, "examplekey123456")
	require.NoError(t, err)
	require.NotEmpty(t, cipherText)

	decryptedText, err := Decrypt(cipherText, "examplekey123456")
	require.NoError(t, err)
	require.Equal(t, plainText, decryptedText)
}

func TestErrorEncrypt(t *testing.T) {
	_, err := Encrypt("", "")
	require.Error(t, err)
}

func TestErrorDecrypt(t *testing.T) {
	_, err := Decrypt("", "")
	require.Error(t, err)
}
