package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/responses"
)

type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (ping PingController) Ping(ctx *gin.Context) {
	res := responses.NewPingResponse()
	ctx.JSON(http.StatusOK, res)
	return
}
