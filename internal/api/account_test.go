package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/mock"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util/token"
	"github.com/stretchr/testify/require"
)

func TestGetAccountByIdAPI(t *testing.T) {
	// buat user random
	user, _ := randomUser(t)

	// buat akun
	account := randomAccount(user.Username)

	// untuk mendapatkan 100% coverage
	testCases := []struct {
		name          string
		accountId     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockMockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder) // check output untuk API
	}{
		{
			name:      "OK",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockMockStore) {
				// build topik atau stubs untuk mock store
				// Yaitu get account'
				// sebanyak 1 kali
				// ekspektasi dari output adalah objek akun dan nil yang berarti tidak error
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)

				// check response body
				requiredBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			// build topik atau stubs untuk mock store
			// Yaitu get account'
			// sebanyak 1 kali
			// ekspektasi dari output adalah objek akun dan nil yang berarti tidak error
			buildStubs: func(store *mockdb.MockMockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockMockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockMockStore) {
				// build topik atau stubs untuk mock store
				// Yaitu get account'
				// sebanyak 1 kali
				// ekspektasi dari output adalah objek akun dan nil yang berarti tidak error
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockMockStore) {
				// build topik atau stubs untuk mock store
				// Yaitu get account'
				// sebanyak 1 kali
				// ekspektasi dari output adalah objek akun dan nil yang berarti tidak error
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountId: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockMockStore) {
				// build topik atau stubs untuk mock store
				// Yaitu get account'
				// sebanyak 1 kali
				// ekspektasi dari output adalah objek akun dan nil yang berarti tidak error
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		// jalankan
		t.Run(tc.name, func(t *testing.T) {
			// deklarasi controller
			ctrl := gomock.NewController(t)

			// Jika program sudah selesai maka defer ke finsih
			defer ctrl.Finish()

			// buat store baru
			store := mockdb.NewMockMockStore(ctrl)
			tc.buildStubs(store)
			// test http server dan send getaccount request
			// digunakan untuk record response
			// server := NewServer(store)
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			// deklarasi url path dari API
			url := fmt.Sprintf("/accounts/%d", tc.accountId)

			// membuat http request
			request, err := http.NewRequest(http.MethodGet, url, nil)

			// Apabila gagal mengembalikan err
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 5000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requiredBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// baca semua data response dari body
	data, err := io.ReadAll(body)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Deklarasi db.account, digunakan untuk menyimpan data object ke db.account
	var gotAccount db.Account

	// Mengembalikan objek dan akan disimpan ke gotAccount
	err = json.Unmarshal(data, &gotAccount)

	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
