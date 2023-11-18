package util

import (
	"fmt"
	"math/rand"
	"strings"
)

// Dipanggil secara langsung ketika package dipanggil pertama kali karena ada init()
// Go akan secara otomatis memberikan angka acak
// func init(){
// Penerepan nilai seed
// Menggunakan parameter int64 dan unixNano
// 	rand.NewSource(time.Now().UnixNano())
// }

// Fungsi menghasilkan bilangan bulat acak antara min dan max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // Mengembalikan nilai 0 -> max-min
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Fungsi menghasilkan string acak
func RandomString(n int) string {
	// Deklarasi pembuat string baru
	var sb strings.Builder

	// Panjang alphabet = 26
	k := len(alphabet)

	for i := 0; i < n; i++ {

		// rand.Intn digunakan untuk mendapatkan posisi acak dari 0 - (k-1)
		c := alphabet[rand.Intn(k)]

		// Tulis karakter yang ada
		sb.WriteByte(c)
	}

	return sb.String()
}

// Fungsi menghasilkan nama acak pada owner
func RandomOwner() string {
	return RandomString(6)
}

// Fungsi menghasilkan uang secara acak
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// Fungsi menghasilkan currency secara acak
func RandomCurrency() string {
	currencies := []string{EUR, USD}
	n := len(currencies)

	// rand.Intn digunakan untuk mendapatkan posisi acak dari 0 - (n-1)
	return currencies[rand.Intn(n)]
}

// Fungsi menghasilkan email secara acak
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
