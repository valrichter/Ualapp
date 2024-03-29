// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateMoneyRecord(ctx context.Context, arg CreateMoneyRecordParams) (MoneyRecord, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccount(ctx context.Context, id int32) error
	DeleteAllAccounts(ctx context.Context) error
	DeleteAllEntries(ctx context.Context) error
	DeleteAllTransfers(ctx context.Context) error
	DeleteAllUsers(ctx context.Context) error
	DeleteMoneyRecordById(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAccountByAccountNumber(ctx context.Context, accountNumber pgtype.Text) (GetAccountByAccountNumberRow, error)
	GetAccountById(ctx context.Context, id int32) (Account, error)
	GetAccountByUserId(ctx context.Context, userID uuid.UUID) ([]Account, error)
	GetEntryById(ctx context.Context, id int32) (Entry, error)
	GetEntryByUserId(ctx context.Context, accountID int32) ([]Entry, error)
	GetMoneyRecordByReference(ctx context.Context, reference string) (MoneyRecord, error)
	GetMoneyRecordsByStatus(ctx context.Context, status string) ([]MoneyRecord, error)
	GetTransferById(ctx context.Context, id int32) (Transfer, error)
	GetTransferFromAccountId(ctx context.Context, fromAccountID int32) ([]Transfer, error)
	GetTransferToAccountId(ctx context.Context, toAccountID int32) ([]Transfer, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateAccountBalanceNew(ctx context.Context, arg UpdateAccountBalanceNewParams) (Account, error)
	UpdateAccountNumber(ctx context.Context, arg UpdateAccountNumberParams) (Account, error)
	UpdateAccountsBalance(ctx context.Context, arg UpdateAccountsBalanceParams) (Account, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
	UpdateUsername(ctx context.Context, arg UpdateUsernameParams) (User, error)
}

var _ Querier = (*Queries)(nil)
