package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

func TestCreateUser(t *testing.T) {
	passwordLength := util.RandomInt(6, 20)
	hashedPassword := util.RandomPassword(passwordLength)

	arg := db.CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user1, err := testQuery.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user1)

	require.Equal(t, arg.Email, user1.Email)
	require.Equal(t, arg.HashedPassword, user1.HashedPassword)
	require.NotZero(t, user1.CreatedAt)
	require.WithinDuration(t, user1.CreatedAt.Time, time.Now(), 2*time.Second)
	require.WithinDuration(t, user1.UpdatedAt.Time, time.Now(), 2*time.Second)

	user2, err := testQuery.CreateUser(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, user2)
}

