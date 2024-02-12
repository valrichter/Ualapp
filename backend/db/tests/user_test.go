package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

// createRandomUser creates a random user for tests
func createRandomUser(t *testing.T) db.User {
	password := util.RandomPassword(util.RandomInt(6, 20))
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.CreatedAt)
	require.WithinDuration(t, user.CreatedAt.Time, time.Now(), 2*time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, time.Now(), 2*time.Second)

	return user
}

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {

	user1 := createRandomUser(t)

	arg := db.CreateUserParams{
		Email:          user1.Email,
		HashedPassword: user1.HashedPassword,
	}

	user2, err := testQuery.CreateUser(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, user2)
}

// TestUpdateUser
func TestUpdateUser(t *testing.T) {

}
