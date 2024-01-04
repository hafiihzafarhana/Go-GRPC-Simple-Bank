package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/middleware"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util/token"
)

// struct ini akan melayani permintaan http
type Server struct {
	config     util.Config
	store      db.MockStore
	router     *gin.Engine
	tokenMaker token.Maker
}

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// fungsi ini membuat instance server baru dan mengatur route api
func NewServer(store db.MockStore, config util.Config) (*Server, error) {
	// inisialisasi token
	// untuk paseto
	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)

	// untuk jwt
	// tokenMaker, err := token.NewJWTMaker(config.PasetoSymmetricKey)

	if err != nil {
		log.Fatal("token maker fail", err)
		return nil, fmt.Errorf("token maker fail %w", err)
	}

	// membuat new server
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	// untuk mendapatkan engine validator yang digunakan oleh Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// buat registrasi terhadap validator
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	// route
	router := gin.Default()

	// User
	// tambah data user
	router.POST("/users", server.createUser)

	// Auth
	// Login
	router.POST("/login", server.login)

	// implementasi middleware
	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker))

	// Account
	// tambah akun
	authRoutes.POST("/accounts", server.createAccount)
	// ambil akun berdasarkan id
	authRoutes.GET("/accounts/:id", server.getAccountById)
	// ambil akun berdasarkan query
	authRoutes.GET("/accounts", server.listAccounts)

	// Transfer
	// tambah data transfer
	authRoutes.POST("/transfers", server.createTransfer)

	// Token
	// Refresh Token
	authRoutes.POST("/token/refresh-token", server.renewAccessToken)

	// tambah route ke router
	server.router = router
}

// fungsi untuk menjalankan server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
