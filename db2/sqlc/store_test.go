package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB) // create a new store instance

	// Arrange: Prepare test data
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
    fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 5 // run 5 concurrent transactions
	amount := int64(10) // each transaction transfer 10 units

	errs := make(chan error) // channel to collect errors
	results := make(chan TransferTxResult) // channel to collect results

	// arrange:  this for loop will start 5 concurrent goroutines
	// each goroutine will call the TransferTx function and send the result to the channel
	for i := 0; i < n; i++ {
        go func() {
           result, err := store.TransferTx(context.Background(), TransferTxParams{
		     FromAccountID: account1.ID,
			 ToAccountID: account2.ID,
			 Amount: amount,
		   })

		   errs <- err // send error to channel
		   results <- result // send result to channel
		}()
    }
	
	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <- errs // receive error from channel
		require.NoError(t, err)

		result := <- results // receive result from channel
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check to make sure that the transfer exists in the database
		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount) // account1 is the sender
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount) // account2 is the receiver
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check to make sure that the entries exist in the database
		_, err = testQueries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		_, err = testQueries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance // original balance - after transfer balance
		diff2 := toAccount.Balance - account2.Balance // after transfer balance - original balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0) // 10, 20, 30, 40, 50 (each transaction transfer 10 units)
		// the amount of transactions is the number of times the amount is transferred
		// 1 * 10, 2 * 10, 3 * 10, 4 * 10, 5 * 10

		k := int(diff1 / amount) // based on the above, k must be between 1 and n
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k) // k must be unique for each transaction
		existed[k] = true // mark k as true to avoid duplicate
	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance - int64(n) * amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updatedAccount2.Balance)

	// Clean up all accounts from the database (including any from previous test runs)
	// Must delete child records first due to foreign key constraints
	_, err = testDB.Exec("DELETE FROM entries")
	require.NoError(t, err)
	_, err = testDB.Exec("DELETE FROM transfers")
	require.NoError(t, err)
	_, err = testDB.Exec("DELETE FROM accounts")
	require.NoError(t, err)

	verifyNoAccountExists(t)

}