package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// inisialisasi untuk mendapatkan objek dari store
	store := NewStore(testDB)
	
	// Membuat akun baru
	account1 := createRandomAccount(t) // pengirim
	account2 := createRandomAccount(t) // penerima

	fmt.Println("Before Transaction", account1.Balance, account2.Balance)

	n := 10 // menjalankan n konkurensi bersamaan transfer transaction
	amount := int64(10) // nilai transfer dari akun 1 ke akun 2

	// karena memang jika pada go routines ini ada kondisi yang tidak terpenuhi, ada kemungkinan goroutines tidak akan berhenti jika hanya menggunaka  require testify
	// maka menggunakan bantuan channels untuk handling atau verifikasi sebuah kesalahan
	// jika ada kesalahan, maka hasilnya akan dikembalikan ke atas (di main go routines)
	// ==================================================================================
	// Channels sendiri digunakan untuk komunikasi antar go routine, koordinasi, dan pertukaran data
	// Belajar: "channels go routine medium"
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// Buat go routines
	for i:=0; i < n; i++ {
		// Debug deadlock detected
		// membuat dan memformat string berdasarkan format yang ditentukan
		txName := fmt.Sprintf("tx %d", i+1)

		go func(){
			// context.Background() digunakan untuk mengelola informasi tambahan yang berkaitan
			// dengan eksekusi program, terutama dalam situasi konkuren, seperti penggunaan goroutine
			
			// Membuat context baru dengan nama transaksi (txName)
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId: account2.ID,
				Amount: amount,
			})

			// Mengirim error ke error channel
			errs <- err

			// Mengirim result ke results channel
			results <- result
		}()
	}

	// existed[1] = true
    // existed[2] = false
    // existed[3] = true
	existed := make(map[int]bool)

	// check results dan adanya error
	for i := 0; i < n; i++ {
		// Mendapatkan data error dari channel error
		err := <-errs

		// Apabila gagal mengembalikan err
		require.NoError(t, err)

		// Mendapatkan data result dari channel results
		result := <-results

		// Check keseluruhan result yang tidak empty
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer

		// Check keseluruhan transfer untuk tidak empty
		require.NotEmpty(t, transfer)

		// Check kesamaan data akun dengan yang ada di transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		// Id transfer dan created_at tidak boleh 0
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Lihat data transfer
		_, err = store.GetTransfer(context.Background(), transfer.ID)

		// Seharusnya tidak error
		require.NoError(t, err)

		// =====================================================

		// check entry pengirim
		// check entries
		fromEntry := result.FromEntry

		// check data tidak kosong
		require.NotEmpty(t, fromEntry)

		// check kesamaan data
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		
		// Nilai id dan created_at entry bukanlah 0
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// ambil data entry sesuai id
		_, err  = store.GetEntry(context.Background(), fromEntry.ID)

		// Seharusnya tidak error
		require.NoError(t, err)

		// =====================================================

		// check entry penerima
		// check entries
		toEntry := result.ToEntry

		// check data tidak kosong
		require.NotEmpty(t, toEntry)

		// check kesamaan data
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		
		// Nilai id dan created_at entry bukanlah 0
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// ambil data entry sesuai id
		_, err  = store.GetEntry(context.Background(), toEntry.ID)

		// Seharusnya tidak error
		require.NoError(t, err)

		// =====================================================
		// Check account

		// Check akun pengirim
		fromAccount := result.FromAccount
		// Seharusnya data tidak empty
		require.NotEmpty(t, fromAccount)

		// Check kesamaan data
		require.Equal(t, account1.ID, fromAccount.ID)

		// Check akun penerima
		toAccount := result.ToAccount

		// Seharusnya data tidak empty
		require.NotEmpty(t, toAccount)

		// Check kesamaan data
		require.Equal(t, account2.ID, toAccount.ID)

		// =====================================================
		// check balance pada akun

		fmt.Println("Look the balance each transaction", fromAccount.Balance, toAccount.Balance)

		// menghitung perbedaan antara input balance akun 1 dan akun 2
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		fmt.Println(diff1)
		fmt.Println(diff2)
		// perbedaan balance akun penerima dengan pengirim harus sama
		require.Equal(t, diff1, diff2)

		// Dan harus memiliki nilai posisitif
		require.True(t, diff1 > 0)

		// nilai perbedaan pada akun 1 harus bisa dibangi sesuai dengan ammount
		// 10 , 20 , 30, 40 , 50 % 10 = 1 , 2 , 3 , 4, 5
		require.True(t, diff1 % amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

		// Seharusnya existed tidak mengandung k (1,2,3,4,5)
		require.NotContains(t, existed, k)
		existed[k] = true // set existed variable ke true
	}

	// Periksa update akhir pada balance akun 1 (pengirim)
	updateAccount1 , err := testQueries.GetAccount(context.Background(), account1.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Periksa update akhir pada balance akun 2 (penerima)
	updateAccount2 , err := testQueries.GetAccount(context.Background(), account2.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	fmt.Println("After Transaction", updateAccount1.Balance, updateAccount2.Balance)

	// saldo akun 1 harus berkurang sebesar n * amount
	require.Equal(t, account1.Balance - int64(n) * amount, updateAccount1.Balance)

	// saldo akun 2 harus bertambah sebesar n * amount
	require.Equal(t, account2.Balance + int64(n) * amount, updateAccount2.Balance)
}

// Konkurensi seperti ini menyebabkan potensi deadlock, jika menggunakan db transaction
func TestTransferTxDeadLock(t *testing.T) {
	// inisialisasi untuk mendapatkan objek dari store
	store := NewStore(testDB) 
	
	// Membuat akun baru
	account1 := createRandomAccount(t) // pengirim
	account2 := createRandomAccount(t) // penerima

	fmt.Println("Before Transaction ", account1.Balance, account2.Balance)

	n := 10 // menjalankan n konkurensi bersamaan transfer transaction
	amount := int64(10) // nilai transfer dari akun 1 ke akun 2

	// karena memang jika pada go routines ini ada kondisi yang tidak terpenuhi, ada kemungkinan goroutines tidak akan berhenti jika hanya menggunaka  require testify
	// maka menggunakan bantuan channels untuk handling atau verifikasi sebuah kesalahan
	// jika ada kesalahan, maka hasilnya akan dikembalikan ke atas (di main go routines)
	// ==================================================================================
	// Channels sendiri digunakan untuk komunikasi antar go routine, koordinasi, dan pertukaran data
	// Belajar: "channels go routine medium"
	errs := make(chan error)
	// results := make(chan TransferTxResult)

	// Buat go routines
	for i:=0; i < n; i++ {

		fromAccountId := account1.ID
		toAccountId := account2.ID

		// apakah penghitungan ganjil atau genap
		// Digunakan agar saling mengirim
		if i % 2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		// Debug deadlock detected
		// membuat dan memformat string berdasarkan format yang ditentukan
		txName := fmt.Sprintf("tx %d", i+1)

		go func(){
			// context.Background() digunakan untuk mengelola informasi tambahan yang berkaitan
			// dengan eksekusi program, terutama dalam situasi konkuren, seperti penggunaan goroutine
			
			// Membuat context baru dengan nama transaksi (txName)
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: fromAccountId,
				ToAccountId: toAccountId,
				Amount: amount,
			})

			// Mengirim error ke error channel
			errs <- err

			// Mengirim result ke results channel
			// results <- result
		}()
	}

	// existed[1] = true
    // existed[2] = false
    // existed[3] = true
	// existed := make(map[int]bool)

	// check results dan adanya error
	for i := 0; i < n; i++ {
		// Mendapatkan data error dari channel error
		err := <-errs

		// Apabila gagal mengembalikan err
		require.NoError(t, err)

		// Mendapatkan data result dari channel results
		// result := <-results

		// Check keseluruhan result yang tidak empty
		// require.NotEmpty(t, result)

		// Check transfer
		// transfer := result.Transfer

		// Check keseluruhan transfer untuk tidak empty
		// require.NotEmpty(t, transfer)

		// Check kesamaan data akun dengan yang ada di transfer
		// require.Equal(t, account1.ID, transfer.FromAccountID)
		// require.Equal(t, account2.ID, transfer.ToAccountID)
		// require.Equal(t, amount, transfer.Amount)

		// Id transfer dan created_at tidak boleh 0
		// require.NotZero(t, transfer.ID)
		// require.NotZero(t, transfer.CreatedAt)

		// Lihat data transfer
		// _, err = store.GetTransfer(context.Background(), transfer.ID)

		// Seharusnya tidak error
		require.NoError(t, err)

		// =====================================================

		// check entry pengirim
		// check entries
		// fromEntry := result.FromEntry

		// check data tidak kosong
		// require.NotEmpty(t, fromEntry)

		// check kesamaan data
		// require.Equal(t, account1.ID, fromEntry.AccountID)
		// require.Equal(t, -amount, fromEntry.Amount)
		
		// Nilai id dan created_at entry bukanlah 0
		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)

		// ambil data entry sesuai id
		// _, err  = store.GetEntry(context.Background(), fromEntry.ID)

		// Seharusnya tidak error
		require.NoError(t, err)

		// =====================================================

		// check entry penerima
		// check entries
		// toEntry := result.ToEntry

		// check data tidak kosong
		// require.NotEmpty(t, toEntry)

		// check kesamaan data
		// require.Equal(t, account2.ID, toEntry.AccountID)
		// require.Equal(t, amount, toEntry.Amount)
		
		// Nilai id dan created_at entry bukanlah 0
		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)

		// ambil data entry sesuai id
		// _, err  = store.GetEntry(context.Background(), toEntry.ID)

		// Seharusnya tidak error
		require.NoError(t, err)

		// =====================================================
		// Check account

		// Check akun pengirim
		// fromAccount := result.FromAccount
		// Seharusnya data tidak empty
		// require.NotEmpty(t, fromAccount)

		// Check kesamaan data
		// require.Equal(t, account1.ID, fromAccount.ID)

		// Check akun penerima
		// toAccount := result.ToAccount

		// Seharusnya data tidak empty
		// require.NotEmpty(t, toAccount)

		// Check kesamaan data
		// require.Equal(t, account2.ID, toAccount.ID)

		// =====================================================
		// check balance pada akun

		// fmt.Println("Look the balance each transaction", fromAccount.Balance, toAccount.Balance)

		// menghitung perbedaan antara input balance akun 1 dan akun 2
		// diff1 := account1.Balance - fromAccount.Balance
		// diff2 := toAccount.Balance - account2.Balance
		// fmt.Println(diff1)
		// fmt.Println(diff2)
		// perbedaan balance akun penerima dengan pengirim harus sama
		// require.Equal(t, diff1, diff2)

		// Dan harus memiliki nilai posisitif
		// require.True(t, diff1 > 0)

		// nilai perbedaan pada akun 1 harus bisa dibangi sesuai dengan ammount
		// 10 , 20 , 30, 40 , 50 % 10 = 1 , 2 , 3 , 4, 5
		// require.True(t, diff1 % amount == 0)
		// k := int(diff1 / amount)
		// require.True(t, k >= 1 && k <= n)

		// Seharusnya existed tidak mengandung k (1,2,3,4,5)
		// require.NotContains(t, existed, k)
		// existed[k] = true // set existed variable ke true
	}

	// Periksa update akhir pada balance akun 1 (pengirim)
	updateAccount1 , err := testQueries.GetAccount(context.Background(), account1.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Periksa update akhir pada balance akun 2 (penerima)
	updateAccount2 , err := testQueries.GetAccount(context.Background(), account2.ID)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	fmt.Println("After Transaction", updateAccount1.Balance, updateAccount2.Balance)

	// saldo akun 1 harus berkurang sebesar n * amount
	require.Equal(t, account1.Balance, updateAccount1.Balance)

	// saldo akun 2 harus bertambah sebesar n * amount
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}