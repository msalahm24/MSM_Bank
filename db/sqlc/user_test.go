package db

import (
	"context"
	"testing"
	"time"

	"github.com/mahmoud24598salah/MSM_Bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPass, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)
	arg := CreateUserParams{
		Username:   util.RandomOwnerName(),
		HashedPass: hashedPass,
		FullName:   util.RandomOwnerName(),
		Email:      util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPass, user.HashedPass)
	require.Equal(t, arg.FullName, user.FullName)
	require.True(t, user.PassChanged.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPass, user2.HashedPass)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
