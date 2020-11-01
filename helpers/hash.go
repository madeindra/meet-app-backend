package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type HashInterface interface {
	Generate(password string) (string, error)
	Verify(hashedPassword string, plainPassword string) error
}

type HashImplementation struct {
	hash HashInterface
}

func NewHashHelper() *HashImplementation {
	return &HashImplementation{}
}

func (hash *HashImplementation) Generate(password string) (string, error) {
	bytePassword := []byte(password)
	byteHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return "", errors.New("Hash failed")
	}

	hashed := string(byteHash)
	return hashed, nil
}

func (hash *HashImplementation) Verify(hashedPassword string, plainPassword string) error {
	byteHash := []byte(hashedPassword)
	bytePassword := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)

	return err
}
