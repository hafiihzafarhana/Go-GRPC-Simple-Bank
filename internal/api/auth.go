package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/exception"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
)

// input untuk login
type loginUserRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"user_name" binding:"required"`
}

// respon setelah berhasil login
type loginUserResponse struct {
	SessionId             uuid.UUID          `json:"session_id"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  createUserResponse `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	// Deklarasi tipe data
	var req loginUserRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// check apakah ada datanya berdasarkan username
	dataUser, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		// jika user tidak ada
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, exception.ErrorResponse(err))
			return
		}

		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	// periksa password apakah sesuai
	err = util.CheckPassword(req.Password, dataUser.Password)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusUnauthorized, exception.ErrorResponse(err))
		return
	}

	// create access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(dataUser.Username, server.config.PasetoAccessTokenDuration)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	// create refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(dataUser.Username, server.config.PasetoRefreshTokenDuration)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     dataUser.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: createUserResponse{
			FullName: dataUser.FullName,
			Username: dataUser.Username,
			Email:    dataUser.Email,
		},
	})
}

// func (server *Server) refreshingToken(ctx *gin.Context) {}
