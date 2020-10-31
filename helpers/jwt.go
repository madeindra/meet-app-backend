package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/madeindra/meet-app/common"
)

var signingKey string = common.GetBearerKey()

//TODO: Change to 1 hour if refresh token already implemented
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