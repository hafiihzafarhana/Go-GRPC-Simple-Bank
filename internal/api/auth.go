package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
)

// input untuk login
type loginUserRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"user_name" binding:"required"`
}

// respon setelah berhasil login
type loginUserResponse struct {
	AccessToken string             `json:"access_token"`
	User        createUserResponse `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	// Deklarasi tipe data
	var req loginUserRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check apakah ada datanya berdasarkan username
	dataUser, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		// jika user tidak ada
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// periksa password apakah sesuai
	err = util.CheckPassword(req.Password, dataUser.Password)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// create token
	token, err := server.tokenMaker.CreateToken(dataUser.Username, server.config.PasetoAccessTokenDuration)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		AccessToken: token,
		User: createUserResponse{
			FullName: dataUser.FullName,
			Username: dataUser.Username,
			Email:    dataUser.Email,
		},
	})
}
