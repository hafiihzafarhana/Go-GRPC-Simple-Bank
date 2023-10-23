package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // Diberi "_" agar tidak hilang pada oleh formatter
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// deklarasi untuk menyimpan kesleuruh method yang ada di Queries
var testQueries *Queries

// mendeklarasikan variabel global baru
var testDB *sql.DB

func TestMain(m *testing.M){
	var err error

	// Membuka DB
	testDB, err = sql.Open(dbDriver, dbSource)
	
	// Lakukan pengecheckan apabila ada error
	if err != nil {
		log.Fatal("Can't connect to DB ", err)
	}

	// New() sudah didefinisikan oleh SQLC
	testQueries = New(testDB)

	// Jalankan dan memberi tahu test lulus atau gagal
	os.Exit(m.Run())
}