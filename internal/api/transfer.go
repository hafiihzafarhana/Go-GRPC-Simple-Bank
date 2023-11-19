package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/exception"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util/token"
)

// input untuk create transfer
type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// controller create transfer
func (server *Server) createTransfer(ctx *gin.Context) {
	// Deklarasi tipe data
	var req transferRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountId, req.Currency)
	// validasi pengirim
	if !valid {
		err := errors.New("invalid data")
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// payload hasil ekstraksi access token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// periksa apakah akun pengirim sesuai
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("account does not belong to authenticated users")
		ctx.JSON(http.StatusUnauthorized, exception.ErrorResponse(err))
		return
	}

	toAccount, valid := server.validAccount(ctx, req.ToAccountId, req.Currency)
	// dan penerima
	if !valid {
		err := errors.New("invalid data")
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// periksa apakah akun penerima sesuai
	if toAccount.ID != req.ToAccountId {
		err := errors.New("account does not belong to authenticated users")
		ctx.JSON(http.StatusUnauthorized, exception.ErrorResponse(err))
		return
	}

	// argumen yang akan dimasukan ke dalam db
	arg := db.TransferTxParams{
		FromAccountId: fromAccount.ID,
		ToAccountId:   toAccount.ID,
		Amount:        req.Amount,
	}

	// Masukan data ke dalam db
	transfer, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}

// fungsi untuk memeriksa apakah akun tersedia
func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	// ambil akun berdasarkan id
	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		// jika akun tidak ada
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, exception.ErrorResponse(err))
			return account, false
		}

		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return account, false
	}

	// jika tidak error, maka periksa currency
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s -> %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, err)
		return account, false
	}

	return account, true
}
