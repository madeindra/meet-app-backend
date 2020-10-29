package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/responses"
)

type CredentialController struct {
	credential models.CredentialInterface
}

func NewCredentialController(credential models.CredentialInterface) *CredentialController {
	return &CredentialController{credential}
}

func (controller *CredentialController) Create(ctx *gin.Context) {
	data := models.NewCredentialData()
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	credential, err := controller.credential.CreateNewCredential(data)

	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responses.NewCredentialResponse(credential.ID, credential.Email)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *CredentialController) FindOne(ctx *gin.Context) {
	data := models.NewCredentialData()
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	credential := controller.credential.FindOneCredential(data)

	if credential.ID == 0 {
		res := responses.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	res := responses.NewAuthenticatedResponse(credential.ID, credential.Email, "fake token")
	ctx.JSON(http.StatusOK, res)
	return
}
