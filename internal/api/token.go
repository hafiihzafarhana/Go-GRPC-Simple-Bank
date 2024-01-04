package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/exception"
)

// input untuk token
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// respon setelah berhasil login
type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	// Deklarasi tipe data
	var req renewAccessTokenRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// check apakah refresh token ada di dalam db
	_, err := server.store.GetSessionByRefreshToken(ctx, req.RefreshToken)

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

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusUnauthorized, exception.ErrorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)

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

	// periksa apakah session diblock atau tidak
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// periksa apakah data session dengan refesh payload memiliki username yang sama
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorect session username")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// periksa refresh token di db session dan refrsh payload sama atau tidak
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched refresh token")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// pada kondisi tertentu, refresh token harus dipaksa kedaluarsa, sehingga harus diperiksa kembali meskipun sudah ada fungsi Valid() di payload.go
	if time.Now().After(session.ExpiresAt){
		err := fmt.Errorf("expired time refresh token")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// create access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.PasetoAccessTokenDuration)

	if err != nil {
		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	})
}
