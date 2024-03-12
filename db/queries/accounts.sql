-- name: CreateAccount :one
INSERT INTO
    accounts (user_id, balance, currency)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountByUserId :many
SELECT * FROM accounts WHERE user_id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateAccountsBalance :one
UPDATE accounts SET balance = $1 WHERE id = $2 RETURNING *;

-- name: UpdateAccountNumber :one
UPDATE accounts SET account_number = $1 WHERE id = $2 RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET
    balance = balance + sqlc.arg (amount)
WHERE
    id = sqlc.arg (id) RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- name: DeleteAllAccounts :exec
DELETE FROM accounts;