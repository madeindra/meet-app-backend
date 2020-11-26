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
	// bind request
	req := entities.NewTokenRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// get email from refresh token
	email, err := controller.bearer.ParseRefresh(req.RefreshToken)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find credential by email in db
	user := controller.credential.New()
	user.Email = email
	userData := controller.credential.FindOne(user)
	if userData.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// find refresh data by user id & refresh token in db
	token := controller.token.New()
	token.UserID = userData.ID
	token.RefreshToken = req.RefreshToken
	data := controller.token.FindOne(token)
	if data.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// generate a new jwt auth token
	authToken, err := controller.bearer.GenerateToken(userData.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// generate a new refresh token
	refreshToken, err := controller.bearer.GenerateRefresh(userData.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// update existing refresh token in db
	data.RefreshToken = refreshToken
	if _, err := controller.token.UpdateByUser(data); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return resposne
	res := entities.NewTokenResponse(authToken, data.RefreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}
