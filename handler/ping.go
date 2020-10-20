package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	res := gin.H{"success": true, "message": "Server is working properly"}
	ctx.JSON(http.StatusOK, res)
	return
}
