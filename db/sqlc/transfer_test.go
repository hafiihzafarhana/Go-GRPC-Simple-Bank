package db

import (
	"context"
	"testing"
	"time"

	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {

	// Create parameter
	arg := CreateTransferParams{
		Amount:        util.RandomMoney(),
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	}

	// Eksekusi create transfer sesuai parameter
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil harusnya mengembalikan transfer
	require.NotEmpty(t, transfer)

	// Memeriksa apakah datanya cocok
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	// Memeriksa apakah id dan timpestamp bukan 0
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	// buat akun
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// tambahkan random transfer berdasarkan akun yang telah dibuat
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	// buat akun
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// tambahkan random transfer berdasarkan akun yang telah dibuat
	transfer1 := createRandomTransfer(t, account1, account2)

	// Mengambil berdasarkan id
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan transfer2
	require.NotEmpty(t, transfer2)

	// Memeriksa apakah datanya cocok
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt.Time, transfer2.CreatedAt.Time, time.Second)
}

func TestListTransfers(t *testing.T) {
	// buat akun
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Buat beberapa akun dan juga transfer
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
	}

	// Deklarasi paramter pada list account
	arg := ListTransfersParams{
		Limit:         3,
		Offset:        4,
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	}

	// List Transfer
	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Panjang harus sesuai dengan limit
	require.Len(t, transfers, 3)

	// Check keseluruhan trasnfers yang telah ditampilkan
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
