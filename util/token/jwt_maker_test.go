package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	// Membuat JWT Maker instance
	maker, err := NewJWTMaker(util.RandomString(32))

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	// Membuat token
	createToken, err := maker.CreateToken(username, duration)

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	// token tidak empty
	require.NotEmpty(t, createToken)

	// Verify token
	payload, err := maker.VerifyToken(createToken)

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	// payload tidak empty
	require.NotEmpty(t, payload)

	// payload id bukanlah null
	require.NotZero(t, payload.ID)

	// cocokan
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	// Membuat JWT Maker instance
	maker, err := NewJWTMaker(util.RandomString(32))

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	username := util.RandomOwner()

	// Membuat token yang expired
	createToken, err := maker.CreateToken(username, -time.Minute)

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	// token tidak empty
	require.NotEmpty(t, createToken)

	// Verify token
	payload, err := maker.VerifyToken(createToken)

	// Harus gagal
	require.Error(t, err)

	// Harus error token expired
	require.EqualError(t, err, ErrTokenExpired.Error())

	// payload harus lah nil
	require.Nil(t, payload)
}

// Tes untuk memeriksa apakah JWT itu invalid
// ini adalah testing untuk memperhatikan alg di header jwt. karena bisa diganti, dan rentan
func TestInvalidJWTTokenAlgNone(t *testing.T) {
	// membuat payload
	payload, err := NewPayload(util.RandomOwner(), time.Minute)

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	// masukan apa saja untuk proses membuat token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	// buat token
	// hanya digunakan pada testing saja
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	// Membuat JWT Maker instance
	maker, err := NewJWTMaker(util.RandomString(32))

	// Apabila gagal, mengembalikan err
	require.NoError(t, err)

	// Verify token
	payload, err = maker.VerifyToken(token)

	// Harus gagal
	require.Error(t, err)

	// Harus error token expired
	require.EqualError(t, err, ErrInvalidToken.Error())

	// payload harus lah nil
	require.Nil(t, payload)
}
