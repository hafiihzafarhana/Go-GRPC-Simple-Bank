package main

import (
	"database/sql"
	"log"

	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/internal/api"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	_ "github.com/lib/pq"
)

func main() {
	// Ambil fungsi load config dalam util
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Can't connect configuration", err)
	}

	// Membuka DB
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	// Lakukan pengecheckan apabila ada error
	if err != nil {
		log.Fatal("Can't connect to DB ", err)
	}

	// Ambil new store
	store := db.NewStore(conn)

	// ambil server
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal("Can't create server", err)
	}

	// pada pengambilan server, maka lanjutkan dengan start
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Can't start server ", err)
	}
}
