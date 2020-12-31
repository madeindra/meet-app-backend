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
	skill      models.SkillInterface
}

func NewProfileController(profile models.ProfilesInterface, credential models.CredentialInterface, skill models.SkillInterface) *ProfilesController {
	return &ProfilesController{profile, credential, skill}
}

func (controller *ProfilesController) GetSingle(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find profile by user id
	data := controller.profile.New()
	data.ID = id

	profile := controller.profile.FindOne(data)
	if profile.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return response
	res := entities.NewProfileResponse(profile.ID, profile.Name, profile.Description, profile.Gender, profile.Picture, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) GetCollections(ctx *gin.Context) {
	// get skillName from url
	skillName := ctx.DefaultQuery("skillName", "")
	userName := ctx.DefaultQuery("userName", "")

	var skillList []uint64

	if skillName != "" {
		// find all skill name stored in db
		userSkill := controller.skill.New()
		userSkill.SkillName = skillName

		skills := controller.skill.FindBy(userSkill)
		for i := range skills {
			skillList = append(skillList, skills[i].UserID)
		}
	}

	if userName != "" {
		// find all user name stored in db
		user := controller.profile.New()
		user.Name = userName

		users := controller.profile.FindBy(user)
		if len(users) == 0 {
			res := entities.NotFoundResponse()
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		// return response
		res := entities.NewProfileBatchResponse(users)
		ctx.JSON(http.StatusOK, res)
		return
	}

	if len(skillList) > 0 {
		profile := controller.profile.FindIn(skillList)
		if len(profile) == 0 {
			res := entities.NotFoundResponse()
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		// return response
		res := entities.NewProfileBatchResponse(profile)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// find all profile stored in db
	profile := controller.profile.FindAll()
	if len(profile) == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return response
	res := entities.NewProfileBatchResponse(profile)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Post(ctx *gin.Context) {
	// bind request
	req := entities.NewProfileRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find credential by user id in db
	userExist := controller.credential.New()
	userExist.ID = req.ID

	if exist := controller.credential.FindOne(userExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	profileExist := controller.profile.New()
	profileExist.ID = req.ID

	// if user already exist, abort process
	if duplicate := controller.profile.FindOne(profileExist); duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// create profile data to insert in db
	data := controller.profile.New()
	data.ID = req.ID
	data.Name = req.Name
	data.Description = req.Description
	data.Gender = req.Gender
	data.Picture = req.Picture
	data.Latitude = req.Latitude
	data.Longitude = req.Longitude

	// insert profile data to db
	profile, err := controller.profile.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewProfileResponse(profile.ID, profile.Name, profile.Description, profile.Gender, profile.Picture, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *ProfilesController) Put(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find profile by user id in db
	checkExisting := controller.profile.New()
	checkExisting.ID = id

	if exist := controller.profile.FindOne(checkExisting); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// bind request
	req := entities.NewProfileRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// create profile data to be updated in db
	data := controller.profile.New()
	data.ID = id
	data.Name = req.Name
	data.Description = req.Description
	data.Gender = req.Gender
	data.Picture = req.Picture
	data.Latitude = req.Latitude
	data.Longitude = req.Longitude

	// update profile data
	profile, err := controller.profile.UpdateByID(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewProfileResponse(profile.ID, profile.Name, profile.Description, profile.Gender, profile.Picture, profile.Latitude, profile.Longitude)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ProfilesController) Delete(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find profile by used id in db
	data := controller.profile.New()
	data.ID = id

	if profile := controller.profile.FindOne(data); profile.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// delete profile data from db
	if err := controller.profile.Delete(data); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return empty response
	ctx.JSON(http.StatusNoContent, nil)
	return
}
