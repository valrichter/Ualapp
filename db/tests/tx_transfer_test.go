package db_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

func TestTransferTx(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	account1 := createRandomAccount(t, user1)
	account2 := createRandomAccount(t, user2)

	fmt.Println(">> Before [", "from:", account1.Balance, "to:", account2.Balance, "]")

	amount := int64(10)
	arg := db.TransferTxRequest{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	errs := make(chan error)
	results := make(chan db.TransferTxResponse)

	n := util.RandomInt(3, 10)
	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			tx, err := testStore.TransferTx(context.Background(), arg)

			errs <- err
			results <- tx
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		tx := <-results
		require.NotEmpty(t, tx)

		// check transfer
		transfer := tx.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check entries
		fromEntry := tx.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := tx.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromAccount := tx.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account1.UserID, fromAccount.UserID)

		toAccount := tx.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		require.Equal(t, account2.UserID, toAccount.UserID)

		// check balances
		fmt.Println(">> tx: [", "from:", fromAccount.Balance, "to:", toAccount.Balance, "]")

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amoun

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> After: [", "from:", updatedAccount1.Balance, "to:", updatedAccount2.Balance, "]")

	totalTransfered := int64(n) * amount
	require.Equal(t, account1.Balance-totalTransfered, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+totalTransfered, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	account1 := createRandomAccount(t, user1)
	account2 := createRandomAccount(t, user2)

	fmt.Println(">> Before [", "from:", account1.Balance, "to:", account2.Balance, "]")

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), db.TransferTxRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> After: [", "from:", updatedAccount1.Balance, "to:", updatedAccount2.Balance, "]")
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
