package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signingKey = "signingkey"
)

func CreateJWT(email string) (string, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(token string) (*jwt.Token, error) {
	validated, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	return validated, err
}

func ParseJWT(token string) (jwt.MapClaims, error) {
	validated, err := ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	if claims, ok := validated.Claims.(jwt.MapClaims); ok && validated.Valid {
		return claims, nil
	}

	return nil, errors.New("Unable to parse token")
}
