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
SELECT * FROM accounts WHERE id = $1 LIMIT 1;
-- Here GetAccount is the name of the function in generated go code :one means one row

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;
-- LIMIT $1 enable pagination so that we only display certain number of rows

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2 WHERE id = $1
RETURNING *;
-- Here UpdateAccount is the name of the function in generated go code :one means return one row

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;