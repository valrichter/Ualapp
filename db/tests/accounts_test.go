package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

func createRandomAccount(t *testing.T, user db.User) db.Account {
	arg := db.CreateAccountParams{
		UserID:   user.ID,
		Balance:  util.RandomMoney(0, 1000),
		Currency: "ARS",
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)
	require.Equal(t, arg.UserID, account.UserID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.CreatedAt)
	require.WithinDuration(t, account.CreatedAt, time.Now(), 2*time.Second)

	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	createRandomAccount(t, user)
}

func TestGetAccountById(t *testing.T) {
	user := createRandomUser(t)
	account1 := createRandomAccount(t, user)

	account2, err := testStore.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.UserID, account2.UserID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 2*time.Second)
}

func TestGetAccountsFromUserId(t *testing.T) {
	user := createRandomUser(t)
	var accounts []db.Account
	amountAccounts := 1

	for i := 0; i < amountAccounts; i++ {
		accounts = append(accounts, createRandomAccount(t, user))
	}

	accountsFromDB, err := testStore.GetAccountsFromUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountsFromDB)
	require.Len(t, accountsFromDB, len(accounts))
}
