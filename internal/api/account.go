package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/exception"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util/token"
	"github.com/lib/pq"
)

// input untuk create account
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// controller create account
func (server *Server) createAccount(ctx *gin.Context) {
	// Deklarasi tipe data
	var req createAccountRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// payload hasil ekstraksi access token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// argumen yang akan dimasukan ke dalam db
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	// Masukan data ke dalam db
	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		// Memeriksa kesalahan 403 (Forbidden)
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, exception.ErrorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// controller get account by id
func (server *Server) getAccountById(ctx *gin.Context) {
	var req getAccountRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindUri(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// argumen yang akan dimasukan ke dalam db
	user_id := req.ID

	// Cari akun
	account, err := server.store.GetAccount(ctx, user_id)

	if err != nil {
		// jika akun tidak ada
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, exception.ErrorResponse(err))
			return
		}

		// jika ada error pada server
		ctx.JSON(http.StatusInternalServerError, exception.ErrorResponse(err))
		return
	}

	// payload hasil ekstraksi access token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := errors.New("account does not belong to authenticated users")
		ctx.JSON(http.StatusUnauthorized, exception.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=5,max=10"`
}

// controller get all account
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest

	// periksa jika req data tidak sesuai
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// Kembalikan response error
		ctx.JSON(http.StatusBadRequest, exception.ErrorResponse(err))
		return
	}

	// payload hasil ekstraksi access token
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// argumen yang akan dimasukan ke dalam db
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	// Cari akun
	account, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, exception.ErrorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, account)
}
