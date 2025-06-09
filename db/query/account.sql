-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner=$1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListAccountsAll :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;



-- name: UpdateAccountData :one
UPDATE accounts
SET
  owner = sqlc.arg(owner),
  currency = sqlc.arg(currency)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateAccountAndGet :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNing *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;