// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: entries.sql

package db

import (
	"context"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO
    entries (account_id, amount)
VALUES ($1, $2) RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int32 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRow(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAllEntries = `-- name: DeleteAllEntries :exec
DELETE FROM accounts
`

func (q *Queries) DeleteAllEntries(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deleteAllEntries)
	return err
}

const getEntryById = `-- name: GetEntryById :one
SELECT id, account_id, amount, created_at FROM entries WHERE id = $1
`

func (q *Queries) GetEntryById(ctx context.Context, id int32) (Entry, error) {
	row := q.db.QueryRow(ctx, getEntryById, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getEntryByUserId = `-- name: GetEntryByUserId :many
SELECT id, account_id, amount, created_at FROM entries WHERE account_id = $1
`

func (q *Queries) GetEntryByUserId(ctx context.Context, accountID int32) ([]Entry, error) {
	rows, err := q.db.Query(ctx, getEntryByUserId, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
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

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries ORDER BY id LIMIT $1 OFFSET $2
`

type ListEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.Query(ctx, listEntries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
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
