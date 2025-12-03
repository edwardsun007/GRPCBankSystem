-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *; -- the * means return all the columns

-- name: GetAccount :one
-- There was a bug here during concurrent transactions:
-- when we read the row, we need to lock it from begin to the commit of the transaction
-- otherwise, the other transaction will get incorrect info because the previous transaction is not completed yet
SELECT * FROM accounts WHERE id = $1 LIMIT 1;
-- Here GetAccount is the name of the function in generated go code :one means one row

-- name: GetAccountForUpdate :one
-- we should use this function instead for transaction
SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;
-- LIMIT $1 enable pagination so that we only display certain number of rows

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2 WHERE id = $1
RETURNING *;
-- Here UpdateAccount is the name of the function in generated go code :one means return one row

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;