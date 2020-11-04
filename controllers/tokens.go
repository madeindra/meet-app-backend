package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/responses"
)

type TokenController struct {
	token      models.TokenInterface
	credential models.CredentialInterface
	bearer     helpers.JWTInterface
}

func NewTokenController(token models.TokenInterface, credential models.CredentialInterface, bearer helpers.JWTInterface) *TokenController {
	return &TokenController{token, credential, bearer}
}

func (controller *TokenController) Refresh(ctx *gin.Context) {
	data := controller.token.New(0, "")
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	email, err := controller.bearer.ParseRefresh(data.RefreshToken)
	if err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user := controller.credential.New(email, "")
	userData := controller.credential.FindOne(user)
	if userData.ID == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	tokenData := controller.token.New(userData.ID, data.RefreshToken)
	token := controller.token.FindByUser(tokenData)
	if token.ID == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	authToken, err := controller.bearer.GenerateToken(userData.Email)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshToken, err := controller.bearer.GenerateRefresh(userData.Email)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	token.RefreshToken = refreshToken
	if _, err := controller.token.Update(token); err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responses.NewTokenResponse(authToken, token.RefreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}
