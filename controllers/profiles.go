package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/responses"
)

type ProfilesController struct {
	profile models.ProfilesInterface
}

func NewProfileController(profile models.ProfilesInterface) *ProfilesController {
	return &ProfilesController{profile}
}

func (controller *ProfilesController) GetSingle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := controller.profile.New(id, "", "", "", 0, 0)
	res := controller.profile.FindByUser(data)
	if res.ID == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) GetCollections(ctx *gin.Context) {
	res := controller.profile.FindAll()
	if len(res) == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Post(ctx *gin.Context) {
	data := controller.profile.New(0, "", "", "", 0, 0)
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := controller.profile.Create(data)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *ProfilesController) Put(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := controller.profile.New(id, "", "", "", 0, 0)

	res := controller.profile.FindByUser(data)
	if res.ID == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Disable update ID
	res, err = controller.profile.UpdateByUser(data)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := controller.profile.New(id, "", "", "", 0, 0)
	res := controller.profile.FindByUser(data)
	if res.ID == 0 {
		res := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if err := controller.profile.DeleteByUser(data); err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
	return
}
