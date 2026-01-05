package api

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	mockdb "github.com/techschool/simple-bank/db2/mock"
	db "github.com/techschool/simple-bank/db2/sqlc"
	"github.com/techschool/simple-bank/utils"
	"go.uber.org/mock/gomock"
	"io"
	"bytes"
	"encoding/json"
)
func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	
	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
			    // i expect the GetAccount method to be called with any context and this account.id exactly
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
				// the return of this stubs should match the return type of the GetAccount method
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// TODO: add more cases later
	}


	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
	
			store := mockdb.NewMockStore(ctrl)
			// since we defined the buildStubs function in the testcase, here we call it with the store
			tc.buildStubs(store)
			
			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
	
			url := fmt.Sprintf("/accounts/%d", account.ID)
			fmt.Println("url:", url)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
	
			// this sends the request through the router and records its response in the recorder itself
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:     utils.RandomInt(1, 1000),
		Owner:  utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

// verify the response body matches the expected account data
// parameter account is the expected account data
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account){
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}