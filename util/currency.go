package util

// currency yang support
const (
	USD = "USD"
	EUR = "EUR"
)

// Fungsi yang digunakan untuk memeriksa apakah currency support atau tidak
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR:
		return true
	default:
		return false
	}
}