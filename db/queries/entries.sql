-- name: CreateEntry :one
INSERT INTO
    entries (account_id, amount)
VALUES ($1, $2) RETURNING *;

-- name: GetEntryById :one
SELECT * FROM entries WHERE id = $1;

-- name: GetEntryByUserId :many
SELECT * FROM entries WHERE account_id = $1;

-- name: ListEntries :many
SELECT * FROM entries ORDER BY id LIMIT $1 OFFSET $2;

-- name: DeleteAllEntries :exec
DELETE FROM accounts;