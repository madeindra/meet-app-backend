package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type HashInterface interface {
	GenerateHash(password string) (string, error)
	VerifyHash(hashedPassword string, plainPassword string) error
}

type HashImplementation struct {
	hash HashInterface
}

func NewHashImplementation() *HashImplementation {
	return &HashImplementation{}
}

func (hash *HashImplementation) GenerateHash(password string) (string, error) {
	bytePassword := []byte(password)
	byteHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return "", errors.New("Hash failed")
	}

	hashed := string(byteHash)
	return hashed, nil
}

func (hash *HashImplementation) VerifyHash(hashedPassword string, plainPassword string) error {
	byteHash := []byte(hashedPassword)
	bytePassword := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)

	return err
}
