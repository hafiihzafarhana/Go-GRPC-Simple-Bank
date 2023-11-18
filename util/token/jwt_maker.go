package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// Ini adalah struct untuk pengembangan JWT
type JWTMaker struct {
	secretKey string
}

// CreateToken implements Maker.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// Membuat payload
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	// Membuat token baru
	// var jwt.SigningMethodHS256 *jwt.SigningMethodHMAC
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// mengembalikan jwt token dalam bentuk string
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken implements Maker.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// sebuah fungsi untuk menentukan metode verifikasi token
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		// Ini tujuanya agar peretas, pada saat mengganti header alg
		// token akan menjadi invalid
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyfunc)

	if err != nil {
		// Karena JWT sebenarnya sudah mengembalikan nilai errornya di file parser.go
		verr, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(verr.Inner, ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

// fungsi yang mengembalikan instance JWTMaker
/* Dalam Go, sebuah tipe yang mengimplementasikan suatu
interface secara otomatis dianggap sebagai instance dari interface tersebut. */
func NewJWTMaker(secretKey string) (Maker, error) {
	// Secret key panjangnya harus sesuai
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secret Key should > Minimum Secret Key")
	}

	// mengambil blue print dari interface maker
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}
