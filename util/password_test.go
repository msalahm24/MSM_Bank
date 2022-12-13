package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := RandomString(8)
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)
	err = CheckPass(hashed, password)
	require.NoError(t, err)

	wrongPass := RandomString(8)
	err = CheckPass(hashed, wrongPass)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
