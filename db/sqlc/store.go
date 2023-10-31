package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Digunakan untuk menggunakan DB tiruan
type MockStore interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store menyediakan semua fungsi untuk menjalankan kueri database satu per satu, serta kombinasi dalam 1 transaksi
// Single query berasal dari SQLC
type Store struct {
	// Harus menyematkan Queries di dalam store untuk memperluas jangkauan fungsi yang hanya menerapkan 1 tabel saja
	// Menggunakan bintang * agar Store memiliki referensi dari Queries, bukan hanya salinanya saja
	// Dan juga apabila ada perubahan bagian kode, akan langsung tercermin di tempat lain
	*Queries // Disebut komposisi atau compotition, memperluas fungsi yang dimiliki Queries, dan bukan inherintance

	// untuk membuat transaksi baru
	db *sql.DB 
}

// karena mockstore itu interface, maka return nya bisa apa saja
func NewStore(db *sql.DB) MockStore {
	return &Store{
		db: db, // mengembalikan objek miliki Store
		Queries: New(db), // mengembalikan objek query
	}
}

// Menjalankan transaksi database generik
// dibutuhkan konteks (untuk transaksi database. 
// lalu konteks digunakan untuk mengendalikan waktu tunggu, pembatalan, dan informasi lain yang berkaitan dengan eksekusi transaksi.) 
// dan callback function
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Memulai db transaksi
	// akan terus berjalan apabila sudah di commited (jika berhasil) atau rollback (jika gagal)
	// TxOptions{} digunakan untuk mengatur isolasi kustom level untuk transaksi
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	// Memanggil fungsi transaksi di sesuai Queries
	q := New(tx)

	// memnaggil fungsi yang ada di queries
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbErr)
		}

		return err
	}

	// lakukan commit
	return tx.Commit()
}

// Berisi input parameter untuk transfer antara 2 akun
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

// Berisi output untuk transfer antara 2 akun
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

// mendeklarasikan variabel txKey untuk membantu mendapatkan data nama transaksi
// membuat objek baru dari struct
var txKey = struct{}{}

func addMoney(ctx context.Context, txName interface{}, q *Queries, accountId1, amount1, accountId2, amount2 int64) (account1, account2 Account, err error){
	// tambah balance pada akun 1
	fmt.Println(txName, "Remove Account Balance From ", account1.ID)
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId1,
		Amount: amount1,
	})

	if err != nil {
		// ini sama dengan return account1, account2, err
		return
	}

	// tambah balance pada akun 2
	fmt.Println(txName, "Add Account Balance to ", account2.ID)
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId2,
		Amount: amount2,
	})

	return
}

// Untuk melakukan transaksi pengiriman uang
// Membuat catatan transfer baru, tambahkan entri akun baru, dan memperbaharui akun balance dalam 1 transaksi db (db transaction)
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// memnaggil fungsi execTx
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Konteks memegang nama transaksi
		txName := ctx.Value(txKey) // nama transaksi
		
		fmt.Println(txName, "Create Transfer")

		// Lakukan penambahan untuk data transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID: arg.ToAccountId,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "Create Entry 1")

		// Menambah entry dari orang pengirim
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			Amount: -arg.Amount,
			AccountID: arg.FromAccountId,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "Create Entry 2")

		// Menambah entry dari orang penerima
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			Amount: arg.Amount,
			AccountID: arg.ToAccountId,
		})

		if err != nil {
			return err
		}

		// dahulukan akun dengan id terkecil terlebih dahulu
		// jika tidak, maka dilakukan sebaliknya
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, txName, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		
			if err != nil {
				return err
			}
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, txName, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err
}
