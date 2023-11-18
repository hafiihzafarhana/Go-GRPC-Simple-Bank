package db

import (
	"context"
	"testing"
	"time"

	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {

	// Create parameter
	arg := CreateEntryParams{
		Amount:    util.RandomMoney(),
		AccountID: account.ID,
	}

	// Eksekusi create entry sesuai parameter
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil harusnya mengembalikan entry
	require.NotEmpty(t, entry)

	// Memeriksa apakah datanya cocok
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)

	// Memeriksa apakah id dan timpestamp bukan 0
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	// buat akun
	account := createRandomAccount(t)

	// tambahkan random entry berdasarkan akun yang telah dibuat
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	// Membuat akun terlebih dahulu dan langsung membuat entry nya
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	// Mengambil berdasarkan id
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan entry2
	require.NotEmpty(t, entry2)

	// Memeriksa apakah datanya cocok
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, time.Second)
}

func TestListEntries(t *testing.T) {
	// Membuat akun terlebih dahulu
	account := createRandomAccount(t)

	// Buat beberapa akun dan juga entry
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	// Deklarasi paramter pada list account
	arg := ListEntriesParams{
		Limit:     2,
		Offset:    3,
		AccountID: account.ID,
	}

	// List Entries
	entries, err := testQueries.ListEntries(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Panjang harus sesuai dengan limit
	require.Len(t, entries, 2)

	// Check keseluruhan entries yang telah ditampilkan
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
