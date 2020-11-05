package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/entities"
)

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getBearerToken(ctx)
		if err != nil {
			res := entities.UnauthorizedResponse()
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		validated, err := validateBearerToken(token)
		if err != nil {
			res := entities.UnauthorizedResponse()
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if claims, ok := validated.Claims.(jwt.MapClaims); ok && validated.Valid {
			ctx.Set("sub", claims["sub"])
		}

		ctx.Next()
	}
}

func getBearerToken(ctx *gin.Context) (string, error) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("Invalid Header")
	}

	extractedToken := strings.Split(token, "Bearer ")
	if len(extractedToken) != 2 {
		return "", errors.New("Invalid token")
	}

	return extractedToken[1], nil
}

func validateBearerToken(token string) (*jwt.Token, error) {
	signingKey := common.GetBearerKey()

	validated, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	return validated, err
}
