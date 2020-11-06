package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
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
	data := entities.NewTokenRequest()
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	email, err := controller.bearer.ParseRefresh(data.RefreshToken)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user := controller.credential.New()
	user.Email = email
	userData := controller.credential.FindOne(user)
	if userData.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	tokenData := controller.token.New()
	tokenData.UserID = userData.ID
	tokenData.RefreshToken = data.RefreshToken
	token := controller.token.FindByUser(tokenData)
	if token.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	authToken, err := controller.bearer.GenerateToken(userData.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshToken, err := controller.bearer.GenerateRefresh(userData.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	token.RefreshToken = refreshToken
	if _, err := controller.token.Update(token); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewTokenResponse(authToken, token.RefreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}
