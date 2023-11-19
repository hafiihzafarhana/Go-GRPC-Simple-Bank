package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/exception"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util/token"
)

// header req
const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// penggunaan
// errors.New -> tidak dapat custom, sehingga hanya sesuai dengan yang dituliskan saja
// fmt.Errorf -> agar dapat custom errornya dengan variabel

// ordinary fungsi
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// anonymous function
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// jika tidak ada header
		if len(authorizationHeader) == 0 {
			err := errors.New("header not exist")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.ErrorResponse(err))
			return
		}

		// jika authorization header ada
		fields := strings.Fields(authorizationHeader)

		// tetapi header field berisi kurang dari 2 string
		if len(fields) < 2 {
			err := errors.New("invalid authorization header fromat")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.ErrorResponse(err))
			return
		}

		// ambil tipe
		authorizationType := strings.ToLower(fields[0])

		// atau tipenya tidak valid
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsuppoted authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.ErrorResponse(err))
			return
		}

		// ambil access tokennya
		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.ErrorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
