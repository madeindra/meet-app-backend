package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ValidateBasic(ctx *gin.Context) (string, string, error) {
	if user, password, ok := ctx.Request.BasicAuth(); ok {
		return user, password, nil
	}

	return "", "", errors.New("Invalid token")
}
