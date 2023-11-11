package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	// akses nilai bidang yang divalidasi
	if currency, ok := fl.Field().Interface().(string); ok {
		// periksa currency sudah support atau tidak
		return util.IsSupportedCurrency(currency)
	}

	return false
}