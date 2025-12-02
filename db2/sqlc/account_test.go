package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simple-bank/utils"
)

// every test function must start with "Test" in the go signature
// createRandomAccount creates a random account for testing but IT IS NOT A UNIT TEST
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	// Assert: Verify the results
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func verifyNoAccountExists(t *testing.T) {
	// Count all accounts using the stored database connection
	var count int
	err := testDB.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count)
	t.Logf("Total accounts in database: %d", count)
}

func TestCreateAccount(t *testing.T) {
	// Arrange: Prepare test data
	account := createRandomAccount(t)

	// Verify the account actually exists in the database by querying it
	retrievedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ID, retrievedAccount.ID)

	// Act: clean up the account created for the test
	err = testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	verifyNoAccountExists(t)
}

func TestDeleteAccount(t *testing.T) {
	// First create an account to delete
	account := createRandomAccount(t)

	// Act: Execute the delete function
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	// Count all accounts using the stored database connection
	verifyNoAccountExists(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)

	// clean up the account created for the test
	err = testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	verifyNoAccountExists(t)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)

	err = testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	verifyNoAccountExists(t)
}

// func TestListAccounts(t *testing.T) {
// 	// Store all created account IDs for cleanup
// 	var createdAccountIDs []int64

// 	// Create 10 accounts and store their IDs
// 	for i := 0; i < 10; i++ {
// 		account := createRandomAccount(t)
// 		createdAccountIDs = append(createdAccountIDs, account.ID)
// 	}

// 	arg := ListAccountsParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}
// 	accounts, err := testQueries.ListAccounts(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, accounts)
// 	t.Logf("Total num of accounts: %d", len(accounts))
// 	require.Len(t, accounts, 5)

// 	for _, account := range accounts {
// 		require.NotEmpty(t, account)
// 	}

// 	// Clean up all accounts from the database (including any from previous test runs)
// 	// Must delete child records first due to foreign key constraints
// 	_, err = testDB.Exec("DELETE FROM entries")
// 	require.NoError(t, err)
// 	_, err = testDB.Exec("DELETE FROM transfers")
// 	require.NoError(t, err)
// 	_, err = testDB.Exec("DELETE FROM accounts")
// 	require.NoError(t, err)

// 	verifyNoAccountExists(t)
// }
