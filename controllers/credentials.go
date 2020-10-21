package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/models"
)

func CreateCredential(ctx *gin.Context) {
	var data models.Credentials
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := gin.H{"success": false, "message": "Bad Request"}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	credential, err := models.CreateCredential(data)

	if err != nil {
		res := gin.H{"success": false, "message": "Internal Server Error"}
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := gin.H{"success": true, "data": credential}
	ctx.JSON(http.StatusOK, res)
	return
}
