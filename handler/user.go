package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/model"
)

func UserCreate(ctx *gin.Context) {
	var data model.User
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := gin.H{"success": false, "message": "Bad Request"}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	newUser := model.CreateUser(data)

	res := gin.H{"success": true, "data": newUser}
	ctx.JSON(http.StatusOK, res)
	return
}
