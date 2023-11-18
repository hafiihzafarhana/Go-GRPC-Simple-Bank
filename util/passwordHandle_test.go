package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	// Buat random password
	password := RandomString(6)

	// Hashing password
	hashPass1, err := HashPassword(password)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Tidak boleh kosong
	require.NotEmpty(t, hashPass1)

	// Check apakah compare password berhasil
	err = CheckPassword(password, hashPass1)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Check dengan password yang salah
	wrongPassword := RandomString(6)

	// Check dengan password yang salah
	err = CheckPassword(wrongPassword, hashPass1)

	// Apabila gagal mengembalikan err
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// Hashing password 2 kali
	hashPass2, err := HashPassword(password)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Tidak boleh kosong
	require.NotEmpty(t, hashPass2)

	// Seharusnya ke 2 hash password 1 dan 2 itu berbeda
	require.NotEqual(t, hashPass1, hashPass2)
}
