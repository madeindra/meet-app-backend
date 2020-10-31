package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

//TODO: Create Interface

func GenerateHash(password string) (string, error) {
	bytePassword := []byte(password)
	byteHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return "", errors.New("Hash failed")
	}

	hashed := string(byteHash)
	return hashed, nil
}

func VerifyHash(hashedPassword string, plainPassword string) error {
	byteHash := []byte(hashedPassword)
	bytePassword := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)

	return err
}
