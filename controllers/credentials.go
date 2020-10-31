package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/responses"
)

//TODO: Create refresh token function
type CredentialController struct {
	credential models.CredentialInterface
}

func NewCredentialController(credential models.CredentialInterface) *CredentialController {
	return &CredentialController{credential}
}

func (controller *CredentialController) Register(ctx *gin.Context) {
	data := models.NewCredentialData("", "")
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	hash, err := helpers.GenerateHash(data.Password)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	data.Password = hash

	user := models.NewCredentialData(data.Email, "")
	duplicate := controller.credential.FindOne(user)
	if duplicate.ID != 0 {
		res := responses.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	credential, err := controller.credential.Create(data)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responses.NewCredentialResponse(credential.ID, credential.Email)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *CredentialController) Login(ctx *gin.Context) {
	data := models.NewCredentialData("", "")
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user := models.NewCredentialData(data.Email, "")
	credential := controller.credential.FindOne(user)
	if credential.ID == 0 {
		res := responses.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	err := helpers.VerifyHash(credential.Password, data.Password)
	if err != nil {
		res := responses.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	token, err := helpers.CreateJWT(credential.Email)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshToken, err := helpers.CreateRefreshToken(credential.Email)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responses.NewAuthenticatedResponse(credential.ID, credential.Email, token, refreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}
