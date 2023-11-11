package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
)

// input untuk create transfer
type transferRequest struct {
	FromAccountId    int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountId int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

// controller create transfer
func (server *Server) createTransfer(ctx *gin.Context){
	// Deklarasi tipe data
	var req transferRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// validasi pengirim
	if !server.validAccount(ctx, req.FromAccountId, req.Currency) && !server.validAccount(ctx, req.ToAccountId, req.Currency) {
		return
	}

	// argumen yang akan dimasukan ke dalam db
	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId: req.ToAccountId,
		Amount: req.Amount,
	}

	// Masukan data ke dalam db
	transfer, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}

// fungsi untuk memeriksa apakah akun tersedia
func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	// ambil akun berdasarkan id
	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		// jika akun tidak ada
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	// jika tidak error, maka periksa currency
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s -> %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, err)
		return false
	}

	return true
}