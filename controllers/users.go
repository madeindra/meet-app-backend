package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/models"
)

type UserController struct {
	profile models.ProfilesInterface
}

func NewUserController(profile models.ProfilesInterface) *UserController {
	return &UserController{profile}
}

func (controller *UserController) FindUser(ctx *gin.Context) {
	// get param
	name := ctx.DefaultQuery("name", "")

	// find existing user by name
	profile := controller.profile.New()
	profile.Name = name

	user := controller.profile.FindLike(profile)
	if len(user) == 0 {
		res := entities.NewUserBatchResponse(controller.profile.NewBatch())
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// return response
	res := entities.NewUserBatchResponse(user)
	ctx.JSON(http.StatusOK, res)
	return
}
