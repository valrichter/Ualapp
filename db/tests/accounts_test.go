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
	account := createRandomAccount(t, user)

	accountFromDB, err := testStore.GetAccountById(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)
	require.Equal(t, account.ID, accountFromDB.ID)
	require.Equal(t, account.UserID, accountFromDB.UserID)
	require.Equal(t, account.Balance, accountFromDB.Balance)
	require.Equal(t, account.Currency, accountFromDB.Currency)
	require.WithinDuration(t, account.CreatedAt, accountFromDB.CreatedAt, 2*time.Second)
}

func TestGetAccountsFromUserId(t *testing.T) {
	user := createRandomUser(t)
	var accounts []db.Account
	amountAccounts := 1

	for i := 0; i < amountAccounts; i++ {
		accounts = append(accounts, createRandomAccount(t, user))
	}

	accountFromDB, err := testStore.GetAccountsFromUserId(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)
	require.Len(t, accountFromDB, len(accounts))
}

func TestListAccounts(t *testing.T) {
	var accounts []db.Account
	amountAccounts := 10

	for i := 0; i < amountAccounts; i++ {
		accounts = append(accounts, createRandomAccount(t, createRandomUser(t)))
	}

	arg := db.ListAccountsParams{
		Limit:  int32(amountAccounts),
		Offset: 0,
	}

	accountFromDB, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)
	require.Len(t, accountFromDB, amountAccounts)
	for _, a := range accounts {
		require.Contains(t, accountFromDB, a)
	}
}

// TestUpdateAccountBalance updates an account balance and checks that it was correctly updated
func TestUpdateAccountBalance(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)

	arg := db.UpdateAccountBalanceParams{
		ID:     account.ID,
		Amount: util.RandomMoney(-1000, 1000),
	}

	accountFromDB, err := testStore.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)
	require.Equal(t, account.ID, accountFromDB.ID)
	require.Equal(t, account.UserID, accountFromDB.UserID)
	require.Equal(t, account.Balance+arg.Amount, accountFromDB.Balance)
	require.Equal(t, account.Currency, accountFromDB.Currency)
	require.NotZero(t, account.CreatedAt)
	require.WithinDuration(t, account.CreatedAt, accountFromDB.CreatedAt, 2*time.Second)
}
