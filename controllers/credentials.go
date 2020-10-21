package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/models"
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
		res := gin.H{"success": false, "message": "Bad Request"}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	credential, err := controller.credential.CreateNewCredential(data)

	if err != nil {
		res := gin.H{"success": false, "message": "Internal Server Error"}
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := gin.H{"success": true, "data": credential}
	ctx.JSON(http.StatusOK, res)
	return
}
