package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

func createRandomAccount(t *testing.T) db.Account {
	user := createRandomUser(t)

	arg := db.CreateAccountParams{
		UserID:   user.ID,
		Balance:  util.RandomMoney(),
		Currency: "ARS",
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.UserID, account.UserID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.WithinDuration(t, account.CreatedAt, time.Now(), 2*time.Second)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
