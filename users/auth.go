package users

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	bCryptCost = 12
)

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
	return string(hash), err
}

func passwordMatchesHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
