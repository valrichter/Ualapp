package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	utils "github.com/valrichter/Ualapp/util"
)

func TestCreateUser(t *testing.T) {
	arg := db.CreateUserParams{
		Email:          utils.RandomEmail(),
		HashedPassword: "secret",
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.CreatedAt)
	require.WithinDuration(t, user.CreatedAt.Time, time.Now(), 2*time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, time.Now(), 2*time.Second)
}
