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


-- we should use this function instead for transaction
-- SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR UPDATE; [ this is not ideal ]
-- however, this will create exclusive lock on the row, 
--which means no other transaction can read or write to the row until the current transaction is committed
-- it blocks UPDATE, DELETE and INSERT operations on the row
-- a better way is to use FOR NO KEY UPDATE
-- this only create weaker lock while allow INSERT, while still block modify key column and DELETE
-- name: GetAccountForUpdate :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;


-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;
-- LIMIT $1 enable pagination so that we only display certain number of rows

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2 WHERE id = $1
RETURNING *;
-- Here UpdateAccount is the name of the function in generated go code :one means return one row

-- name: AddAccountBalance :one
-- sqlc.arg(amount) allows use to use the amount variable in generated go code, because balance doesn't make sense
UPDATE accounts SET balance = balance + sqlc.arg(amount) 
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;