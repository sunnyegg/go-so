package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashToken returns a hash of the token using bcrypt
func HashToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing token: %v", err)
	}
	return string(hashedToken), nil
}

func VerifyToken(token, hashedToken string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
}
