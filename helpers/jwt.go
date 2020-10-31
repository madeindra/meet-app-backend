package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/madeindra/meet-app/common"
)

func CreateJWT(email string) (string, error) {
	signingKey := common.GetBearerKey()
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(email string) (string, error) {
	refreshKey := common.GetRefreshKey()
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(refreshKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
