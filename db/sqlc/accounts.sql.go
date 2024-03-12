// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: accounts.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO
    accounts (user_id, balance, currency)
VALUES ($1, $2, $3) RETURNING id, user_id, balance, currency, created_at, account_number
`

type CreateAccountParams struct {
	UserID   int32  `json:"user_id"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount, arg.UserID, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.AccountNumber,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteAccount, id)
	return err
}

const deleteAllAccounts = `-- name: DeleteAllAccounts :exec
DELETE FROM accounts
`

func (q *Queries) DeleteAllAccounts(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteAllAccounts)
	return err
}

const getAccountById = `-- name: GetAccountById :one
SELECT id, user_id, balance, currency, created_at, account_number FROM accounts WHERE id = $1
`

func (q *Queries) GetAccountById(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.AccountNumber,
	)
	return i, err
}

const getAccountByUserId = `-- name: GetAccountByUserId :many
SELECT id, user_id, balance, currency, created_at, account_number FROM accounts WHERE user_id = $1
`

func (q *Queries) GetAccountByUserId(ctx context.Context, userID int32) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAccountByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
			&i.AccountNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, user_id, balance, currency, created_at, account_number FROM accounts ORDER BY id LIMIT $1 OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
			&i.AccountNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccountBalance = `-- name: UpdateAccountBalance :one
UPDATE accounts
SET
    balance = balance + $1
WHERE
    id = $2 RETURNING id, user_id, balance, currency, created_at, account_number
`

type UpdateAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int32 `json:"id"`
}

func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.AccountNumber,
	)
	return i, err
}

const updateAccountsBalance = `-- name: UpdateAccountsBalance :one
UPDATE accounts SET balance = $1 WHERE id = $2 RETURNING id, user_id, balance, currency, created_at, account_number
`

type UpdateAccountsBalanceParams struct {
	Balance int64 `json:"balance"`
	ID      int32 `json:"id"`
}

func (q *Queries) UpdateAccountsBalance(ctx context.Context, arg UpdateAccountsBalanceParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateAccountsBalance, arg.Balance, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.AccountNumber,
	)
	return i, err
}
