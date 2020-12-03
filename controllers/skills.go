package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/models"
)

type SkillsController struct {
	skill   models.SkillInterface
	profile models.ProfilesInterface
}

func NewSkillController(skill models.SkillInterface, profile models.ProfilesInterface) *SkillsController {
	return &SkillsController{skill, profile}
}

func (controller *SkillsController) GetSingle(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find skill by user id
	data := controller.skill.New()
	data.UserID = id

	skill := controller.skill.FindOne(data)
	if skill.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return response
	res := entities.NewSkillResponse(skill.UserID, skill.Name)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *SkillsController) GetCollections(ctx *gin.Context) {
	// find all skill stored in db
	skill := controller.skill.FindAll()
	if len(skill) == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return response
	res := entities.NewSkillBatchResponse(skill)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *SkillsController) Post(ctx *gin.Context) {
	// bind request
	req := entities.NewSkillRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find credential by user id in db
	userExist := controller.profile.New()
	userExist.UserID = req.UserID

	if exist := controller.profile.FindOne(userExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	skillExist := controller.skill.New()
	skillExist.UserID = req.UserID

	// if user already exist, abort process
	if duplicate := controller.skill.FindOne(skillExist); duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// create skill data to insert in db
	data := controller.skill.New()
	data.UserID = req.UserID
	data.Name = req.Name

	// insert skill data to db
	skill, err := controller.skill.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewSkillResponse(skill.UserID, skill.Name)
	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *SkillsController) Put(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find skill by user id in db
	checkExisting := controller.skill.New()
	checkExisting.UserID = id

	if exist := controller.skill.FindOne(checkExisting); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// bind request
	req := entities.NewSkillRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// create skill data to be updated in db
	data := controller.skill.New()
	data.UserID = req.UserID
	data.Name = req.Name

	// update skill data
	skill, err := controller.skill.UpdateByUser(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewSkillResponse(skill.UserID, skill.Name)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *SkillsController) Delete(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find skill by used id in db
	data := controller.skill.New()
	data.UserID = id

	if skill := controller.skill.FindOne(data); skill.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// delete skill data from db
	if err := controller.skill.Delete(data); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return empty response
	ctx.JSON(http.StatusNoContent, nil)
	return
}
