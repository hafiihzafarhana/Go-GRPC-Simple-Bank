package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	_ "github.com/lib/pq" // Diberi "_" agar tidak hilang pada oleh formatter
)

// deklarasi untuk menyimpan kesleuruh method yang ada di Queries
var testQueries *Queries

// mendeklarasikan variabel global baru
var testDB *sql.DB

func TestMain(m *testing.M){
	// Ambil fungsi load config dalam util
	config, err := util.LoadConfig("../../")

	if err != nil {
		log.Fatal("Can't connect configuration", err)
	}

	// Membuka DB
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	
	// Lakukan pengecheckan apabila ada error
	if err != nil {
		log.Fatal("Can't connect to DB ", err)
	}

	// New() sudah didefinisikan oleh SQLC
	testQueries = New(testDB)

	// Jalankan dan memberi tahu test lulus atau gagal
	os.Exit(m.Run())
}