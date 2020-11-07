package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/models"
)

type ProfilesController struct {
	profile    models.ProfilesInterface
	credential models.CredentialInterface
}

func NewProfileController(profile models.ProfilesInterface, credential models.CredentialInterface) *ProfilesController {
	return &ProfilesController{profile, credential}
}

func (controller *ProfilesController) GetSingle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	data := controller.profile.New()
	data.UserID = id

	profile := controller.profile.FindOne(data)
	if profile.ID == 0 {
		profile := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	res := entities.NewProfileResponse(profile.UserID, profile.FirstName, profile.LastName, profile.Description, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) GetCollections(ctx *gin.Context) {
	profile := controller.profile.FindAll()
	if len(profile) == 0 {
		profile := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	res := entities.NewProfileBatchResponse(profile)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Post(ctx *gin.Context) {
	req := entities.NewProfileRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userExist := controller.credential.New()
	userExist.ID = req.UserID

	if exist := controller.credential.FindOne(userExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	profileExist := controller.profile.New()
	profileExist.UserID = req.UserID

	if duplicate := controller.profile.FindOne(profileExist); duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	data := controller.profile.New()
	data.UserID = req.UserID
	data.FirstName = req.FirstName
	data.LastName = req.LastName
	data.Description = req.Description
	data.Latitude = req.Latitude
	data.Longitude = req.Longitude

	profile, err := controller.profile.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewProfileResponse(profile.UserID, profile.FirstName, profile.LastName, profile.Description, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *ProfilesController) Put(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	checkExisting := controller.profile.New()
	checkExisting.UserID = id

	if exist := controller.profile.FindOne(checkExisting); exist.ID == 0 {
		profile := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	req := entities.NewProfileRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		profile := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	data := controller.profile.New()
	data.UserID = req.UserID
	data.FirstName = req.FirstName
	data.LastName = req.LastName
	data.Description = req.Description
	data.Latitude = req.Latitude
	data.Longitude = req.Longitude

	profile, err := controller.profile.UpdateByUser(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewProfileResponse(profile.UserID, profile.FirstName, profile.LastName, profile.Description, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	data := controller.profile.New()
	data.UserID = id

	if profile := controller.profile.FindOne(data); profile.ID == 0 {
		profile := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, profile)
		return
	}

	if err := controller.profile.Delete(data); err != nil {
		profile := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, profile)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
	return
}
