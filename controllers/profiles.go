package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/models"
)

type ProfilesController struct {
	profile models.ProfilesInterface
}

func NewProfileController(profile models.ProfilesInterface) *ProfilesController {
	return &ProfilesController{profile}
}

func (controller *ProfilesController) GetSingle(ctx *gin.Context) {
	return
}

func (controller *ProfilesController) GetCollections(ctx *gin.Context) {
	return
}

func (controller *ProfilesController) Post(ctx *gin.Context) {
	return
}

func (controller *ProfilesController) Put(ctx *gin.Context) {
	return
}

func (controller *ProfilesController) Delete(ctx *gin.Context) {
	return
}
