package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// untuk mendapatkan engine validator yang digunakan oleh Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// buat registrasi terhadap validator
		v.RegisterValidation("currency", validCurrency)
	}

	// Account
	// tambah akun
	router.POST("/accounts", server.createAccount)
	// ambil akun berdasarkan id
	router.GET("/accounts/:id", server.getAccountById)
	// ambil akun berdasarkan query
	router.GET("/accounts", server.listAccounts)

	// Transfer
	// tambah data transfer
	router.POST("/transfers", server.createTransfer)

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
