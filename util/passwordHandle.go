package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hasing sebuah password
func HashPassword(password string) (string, error) {
	// Generate hash sesuai dengan input password sesuai dengan bcrypt.DefaultCost = 10
	// []byte(password) => merubah string ke []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed hash password %w", err)
	}

	return string(hash), err
}

// Memeriksa password
func CheckPassword(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}