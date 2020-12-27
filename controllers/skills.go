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

	// find skill by id
	data := controller.skill.New()
	data.ID = id

	skill := controller.skill.FindOne(data)
	if skill.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return response
	res := entities.NewSkillResponse(skill.ID, skill.UserID, skill.SkillName)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *SkillsController) GetCollections(ctx *gin.Context) {
	// get skillName from url
	skillName := ctx.DefaultQuery("skillName", "")

	// get userId from url
	userID, err := strconv.ParseUint(ctx.DefaultQuery("userId", "0"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find all user's skill or skill name stored in db
	userSkill := controller.skill.New()

	if skillName != "" {
		userSkill.SkillName = skillName
	}

	if userID != 0 {
		userSkill.UserID = userID
	}

	if userSkill.SkillName != "" || userSkill.UserID != 0 {
		skill := controller.skill.FindBy(userSkill)
		if len(skill) == 0 {
			res := entities.NotFoundResponse()
			ctx.JSON(http.StatusNotFound, res)
			return
		}

		res := entities.NewSkillBatchResponse(skill)
		ctx.JSON(http.StatusOK, res)
		return
	}

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
	req := entities.NewSkillBatchRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find credential by user id in db
	userExist := controller.profile.New()
	userExist.ID = req.UserID

	if exist := controller.profile.FindOne(userExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// create batch skill data to insert in db
	bulkData := controller.skill.NewBulk()

	for _, v := range req.SkillName {
		data := controller.skill.New()
		data.UserID = req.UserID
		data.SkillName = v

		// find duplicate user skill
		if duplicate := controller.skill.FindOne(data); duplicate.ID != 0 {
			continue
		}

		bulkData = append(bulkData, data)
	}

	if len(bulkData) == 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// insert bulk skills data to db
	skills, err := controller.skill.Create(bulkData)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewSkillBatchResponse(skills)
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

	// find skill by id in db
	checkExisting := controller.skill.New()
	checkExisting.ID = id

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
	data.ID = id
	data.SkillName = req.SkillName

	// update skill data
	skill, err := controller.skill.UpdateById(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewSkillResponse(skill.ID, skill.UserID, skill.SkillName)
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
	data.ID = id

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
