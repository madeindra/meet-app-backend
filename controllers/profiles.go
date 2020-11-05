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
		profile := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	data := controller.profile.New(id, "", "", "", 0, 0)
	profile := controller.profile.FindByUser(data)
	if profile.ID == 0 {
		profile := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	res := responses.NewProfileResponse(profile.UserID, profile.FirstName, profile.LastName, profile.Description, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) GetCollections(ctx *gin.Context) {
	profile := controller.profile.FindAll()
	if len(profile) == 0 {
		profile := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	res := responses.NewProfileBatchResponse(profile)
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

	checkExisting := controller.profile.New(data.UserID, "", "", "", 0, 0)
	duplicate := controller.profile.FindByUser(checkExisting)
	if duplicate.ID != 0 {
		res := responses.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// TODO: Disable ID Input
	profile, err := controller.profile.Create(data)
	if err != nil {
		res := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := responses.NewProfileResponse(profile.UserID, profile.FirstName, profile.LastName, profile.Description, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *ProfilesController) Put(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	data := controller.profile.New(id, "", "", "", 0, 0)

	profile := controller.profile.FindByUser(data)
	if profile.ID == 0 {
		profile := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		profile := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	// TODO: Disable ID update
	// TODO: Disable userId update
	profile, err = controller.profile.UpdateByUser(data)
	if err != nil {
		profile := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, profile)
		return
	}

	res := responses.NewProfileResponse(profile.UserID, profile.FirstName, profile.LastName, profile.Description, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := responses.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	data := controller.profile.New(id, "", "", "", 0, 0)
	profile := controller.profile.FindByUser(data)
	if profile.ID == 0 {
		profile := responses.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	if err := controller.profile.DeleteByUser(data); err != nil {
		profile := responses.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, profile)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
	return
}
