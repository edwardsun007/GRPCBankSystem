package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all functions to execute db queries and transactions
type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries // Shares the same instance — all Store methods use the same Queries due to pointer
	// Enable direct access, can execute queries like this: store.GetAccount(ctx, 1)
	db *sql.DB
}

// NewStore creates a new Store with the given database connection
// parameter db is the database connection
// returns a new Store instance
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// parameter ctx is the context
// parameter fn is CALLBACK function that will be executed within the transaction
// returns error if the transaction fails
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // nil means use the default transaction options
	if err != nil {
		return err
	}

	q := New(tx) // creates a new Queries instance with the transaction
	// unlike line 18 where a db connection is passed, here a tx is passed
	// it works because New() accepts DBTX type
	err = fn(q)     // execute the callback function passing the new Queries instance
	if err != nil { // if callback function execution fails, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transactionErr: %v, rollbackErr: %v", err, rbErr)
		}
		return err // rollback successful, return the original error
	}
	return tx.Commit() // everything in the transaction succeeded, commit the transaction
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // the created transfer record
	FromAccount Account  `json:"from_account"` // outgoing account after balance is updated
	ToAccount   Account  `json:"to_account"`   // incoming account after balance is updated
	FromEntry   Entry    `json:"from_entry"`   // the created entry record for the outgoing account
	ToEntry     Entry    `json:"to_entry"`     // the created entry record for the incoming account
}

var txKey = struct{}{}

// the seccond bracket means that it is empty object
// struct{} is the type

// TransferTx performs a money transfer from one account to another
// it creates a transfer record, add account entries, and update accounts balance within a single database transaction
// parameter ctx is the context
// parameter arg is the transfer request
// returns error if the transfer fails
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// the second parameter is a callback function that will be executed within the transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var ctError error // declare err here to avoid shadowing the err in the outer scope
		// the next line assign the result of the CreateTransfer to the result.Transfer variable
		// defined in the outer scope
		result.Transfer, ctError = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if ctError != nil {
			return ctError
		}

		// now add the account entries
		var feError error
		result.FromEntry, feError = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if feError != nil {
			return feError // the transaction will be rolled back if this error occurs
		}

		var teError error
		result.ToEntry, teError = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if teError != nil {
			return teError // the transaction will be rolled back if this error occurs
		}

		// steps for updating sender account balance
		if arg.FromAccountID < arg.ToAccountID { // step to avoid deadlock: to fix case where both concurrenttransaction try to update the same account
			var trError error
			result.FromAccount, result.ToAccount, trError = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if trError != nil {
				return trError // the transaction will be rolled back if this error occurs
			}
		} else { // update receiver account balance first
			// steps for updating receiver account balance
			var trError error
			result.ToAccount, result.FromAccount, trError = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if trError != nil {
				return trError // the transaction will be rolled back if this error occurs
			}
		}
		return nil
	}) // this block does the job of creating the transfer record

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64, amount1 int64,
	accountID2 int64, amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
