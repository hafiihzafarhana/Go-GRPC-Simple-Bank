package db

import (
	"context"
	"testing"
	"time"

	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	// Random user
	hashedPassword, err := util.HashPassword(util.RandomString(6))

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// Tidak boleh kosong
	require.NotEmpty(t, hashedPassword)

	// Create parameter
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
	}

	// Eksekusi create user sesuai parameter
	user, err := testQueries.CreateUser(context.Background(), arg)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan user
	require.NotEmpty(t, user)

	// Memeriksa apakah datanya cocok
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)

	// Memeriksa apakah id dan timpestamp bukan 0
	require.NotZero(t, user.CreatedAt)
	// saat pertama kali user dibuat, harapanya bidang ini berisi zero/nol timestampt
	require.True(t, user.PasswordChangeAt.Time.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// Membuat akun terlebih dahulu
	user1 := createRandomUser(t)

	// Mengambil berdasarkan id
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	// Apabila gagal mengembalikan err
	require.NoError(t, err)

	// apabila berhasil mengembalikan user2
	require.NotEmpty(t, user2)

	// Memeriksa apakah datanya cocok
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user1.PasswordChangeAt.Time, user2.PasswordChangeAt.Time, time.Second)
}
