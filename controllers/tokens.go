package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/responses"
)

type TokenController struct {
	token models.TokenInterface
}

func NewTokenController(token models.TokenInterface) *TokenController {
	return &TokenController{token}
}

func (controller *TokenController) GetSingle(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := models.NewTokenData(id, "")

	token := controller.token.FindOne(data)
	if token.ID == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	ctx.JSON(http.StatusOK, token)
	return
}
