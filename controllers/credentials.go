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
	var data models.Credentials
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

	resData := responses.NewCredentialData(credential.ID, credential.Email)
	res := responses.NewCredentialResponse(resData)
	ctx.JSON(http.StatusOK, res)
	return
}
