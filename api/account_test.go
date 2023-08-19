package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/ibnumei/digitalBankGo/db/mock"
	db "github.com/ibnumei/digitalBankGo/db/sqlc"
	"github.com/ibnumei/digitalBankGo/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	store := mockdb.NewMockStore(mockController)
	//build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err  := http.NewRequest(http.MethodGet, url ,nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	//check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 100),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	// gotAccount := db.Account
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
