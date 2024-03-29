package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/madeindra/meet-app/common"
)

type JWTInterface interface {
	GenerateToken(email string) (string, error)
	GenerateRefresh(email string) (string, error)
	ParseRefresh(token string) (string, error)
}

type JWTImplementation struct {
	bearer JWTInterface
}

func NewJWTHelper() *JWTImplementation {
	return &JWTImplementation{}
}

func (bearer *JWTImplementation) GenerateToken(email string) (string, error) {
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

func (bearer *JWTImplementation) GenerateRefresh(email string) (string, error) {
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

func (bearer *JWTImplementation) ParseRefresh(token string) (string, error) {
	refreshKey := common.GetRefreshKey()

	validated, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(refreshKey), nil
	})

	if err != nil {
		return "", errors.New("Failed parsing token")
	}

	if claims, ok := validated.Claims.(jwt.MapClaims); ok && validated.Valid {
		return claims["sub"].(string), nil
	}

	return "", errors.New("Failed parsing token")
}
