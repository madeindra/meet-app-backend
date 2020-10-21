package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (ping PingController) Ping(ctx *gin.Context) {
	res := gin.H{"success": true, "message": "Server is working properly"}
	ctx.JSON(http.StatusOK, res)
	return
}
