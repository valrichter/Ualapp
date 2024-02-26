-- name: CreateTransfer :one
INSERT INTO
    transfers (
        from_account_id, to_account_id, amount
    )
VALUES ($1, $2, $3) RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfers WHERE id = $1;

-- name: GetTransferFromAccountId :many
SELECT * FROM transfers WHERE from_account_id = $1;

-- name: GetTransferToAccountId :many
SELECT * FROM transfers WHERE to_account_id = $1;

-- name: ListTransfers :many
SELECT * FROM transfers ORDER BY id LIMIT $1 OFFSET $2;

-- name: DeleteAllTransfers :exec
DELETE FROM transfers;