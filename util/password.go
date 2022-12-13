package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash the password %s", err)
	}

	return string(hashedPass), nil
}

func CheckPass(hashedpass string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(pass))
}
