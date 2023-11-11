package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/lib/pq"
)

// input untuk create transfer
type createUserRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	FullName    string `json:"full_name"`
	Username string `json:"username"`
	Email string `json:"email"`
}

func (server *Server) createUser(ctx *gin.Context){
	// Deklarasi tipe data
	var req createUserRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPass, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// argumen yang akan dimasukan ke dalam db
	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Password: hashPass,
		Email: req.Email,
	}

	// Masukan data ke dalam db
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		// Memeriksa kesalahan 403 (Forbidden)
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, createUserResponse{
		FullName: user.FullName,
		Username: user.Username,
		Email: user.Email,
	})
}