package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
)

// struct ini akan melayani permintaan http
type Server struct {
	store db.MockStore
	router *gin.Engine
}

// fungsi ini membuat instance server baru dan mengatur route api
func NewServer(store db.MockStore) *Server {
	// membuat new server
	server := &Server{
		store: store,
	}

	// route
	router := gin.Default()

	// tambah akun
	router.POST("/accounts", server.createAccount)
	// ambil akun berdasarkan id
	router.GET("/accounts/:id", server.getAccountById)
	// ambil akun berdasarkan query
	router.GET("/accounts", server.listAccounts)

	// tambah route ke router
	server.router = router
	return server
}

// fungsi untuk menjalankan server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// fungsi untuk mengirimkan error
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
