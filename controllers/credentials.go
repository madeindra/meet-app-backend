package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/responses"
)

type CredentialController struct {
	credential models.CredentialInterface
	token      models.TokenInterface
	hash       helpers.HashInterface
	bearer     helpers.JWTInterface
}

func NewCredentialController(credential models.CredentialInterface, token models.TokenInterface, hash helpers.HashInterface, bearer helpers.JWTInterface) *CredentialController {
	return &CredentialController{credential, token, hash, bearer}
}

func (controller *CredentialController) Register(ctx *gin.Context) {
	data := models.NewCredentialData("", "")
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	hash, err := controller.hash.GenerateHash(data.Password)
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

	err := controller.hash.VerifyHash(credential.Password, data.Password)
	if err != nil {
		res := responses.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	token, err := controller.bearer.CreateJWT(credential.Email)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshToken, err := controller.bearer.CreateRefreshToken(credential.Email)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshTokenData := models.NewTokenData(credential.ID, refreshToken)
	if _, err := controller.token.Update(refreshTokenData); err != nil {
		if _, err := controller.token.Create(refreshTokenData); err != nil {
			res := responses.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	res := responses.NewAuthenticatedResponse(credential.ID, credential.Email, token, refreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}
