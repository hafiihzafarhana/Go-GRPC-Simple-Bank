package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	// Create parameter
	arg := CreateAccountParams{
		Owner: createRandomUser(t).Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// Eksekusi create account sesuai parameter
	account, err := testQueries.CreateAccount(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan account
	require.NotEmpty(t, account)

	// Memeriksa apakah datanya cocok
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// Memeriksa apakah id dan timpestamp bukan 0
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	// Membuat akun terlebih dahulu
	account1 := createRandomAccount(t)

	// Mengambil berdasarkan id
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan account2
	require.NotEmpty(t, account2)

	// Memeriksa apakah datanya cocok
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	// Membuat akun terlebih dahulu
	account1 := createRandomAccount(t)

	// Update Account Params
	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: util.RandomMoney(),
	}
	
	// Update account
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan account2
	require.NotEmpty(t, account2)

	// Memeriksa apakah datanya cocok
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	// Membuat akun terlebih dahulu
	account1 := createRandomAccount(t)

	// Hapus account
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Mengambil berdasarkan id
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	// Apabila ada datanya, berarti merupakan kesalahan
	require.Error(t, err)

	// Apabila ada kesalahan, maka berhasil. harapanya fungsi akan menghasilkan kesalahan
	require.EqualError(t, err, sql.ErrNoRows.Error())

	// Harapanya account2, adalah empy
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	// Buat beberapa akun
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	} 

	// Deklarasi paramter pada list account
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Panjang harus sesuai dengan limit
	require.Len(t, accounts, 5)

	// Check keseluruhan akun yang telah ditampilkan
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}