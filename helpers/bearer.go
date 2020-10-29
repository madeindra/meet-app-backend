package helpers

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBearerToken(ctx *gin.Context) (string, error) {
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
