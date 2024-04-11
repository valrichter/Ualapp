package db

import (
	"context"
)

// Contains the input parameters for creating a new transfer
type TransferTxRequest struct {
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID   int32 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Contains the response of a transfer transaction
type TransferTxResponse struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Money transfer from one account to another.
// It create a transfer record, update account entries, and update account balances in a single database transaction.

func (store *PostgreSQLStore) TransferTx(ctx context.Context, req TransferTxRequest) (TransferTxResponse, error) {

	var res TransferTxResponse

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. Make a transfer
		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(req))
		if err != nil {
			return err
		}

		// 2. Money is moving out of the account 'FromAccountID'
		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: req.FromAccountID,
			Amount:    -req.Amount,
		})
		if err != nil {
			return err
		}

		// 3. Money is moving into the account 'ToAccountID'
		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: req.ToAccountID,
			Amount:    req.Amount,
		})
		if err != nil {
			return err
		}

		// 4. Update accounts with new balances
		// Check FromAccountID < ToAccountID for avoiding deadlock
		if req.FromAccountID < req.ToAccountID {
			res.FromAccount, res.ToAccount, err = UpdateMoney(ctx, q, req.FromAccountID, -req.Amount, req.ToAccountID, +req.Amount)
		} else {
			res.ToAccount, res.FromAccount, err = UpdateMoney(ctx, q, req.ToAccountID, +req.Amount, req.FromAccountID, -req.Amount)
		}

		return err
	})

	return res, err
}

func UpdateMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int32,
	amount1 int64,
	accountID2 int32,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
