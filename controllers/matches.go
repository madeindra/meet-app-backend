package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/models"
)

type MatchController struct {
	match      models.MatchInterface
	credential models.CredentialInterface
}

func NewMatchController(match models.MatchInterface, credential models.CredentialInterface) *MatchController {
	return &MatchController{match, credential}
}

func (controller *MatchController) GetSingle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		match := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, match)
		return
	}

	data := controller.match.New()
	data.ID = id

	match := controller.match.FindOne(data)
	if match.ID == 0 {
		match := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, match)
		return
	}

	res := entities.NewMatchResponse(match.ID, match.UserID, match.UserMatch, match.Liked)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *MatchController) GetCollections(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.DefaultQuery("userId", "0"), 10, 64)
	if err != nil {
		match := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, match)
		return
	}

	userMatch, err := strconv.ParseUint(ctx.DefaultQuery("matchTo", "0"), 10, 64)
	if err != nil {
		match := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, match)
		return
	}

	var boolSensitive bool = false
	var liked bool = false

	if ctx.Query("liked") != "" {
		liked, err = strconv.ParseBool(ctx.Query("liked"))
		if err != nil {
			match := entities.BadRequestResponse()
			ctx.JSON(http.StatusBadRequest, match)
			return
		}

		boolSensitive = true
	}

	data := controller.match.New()
	data.UserID = userID
	data.UserMatch = userMatch
	data.Liked = liked

	match := controller.match.FindBy(data, boolSensitive)
	if len(match) == 0 {
		match := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, match)
		return
	}

	res := entities.NewMatchBatchResponse(match)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *MatchController) Post(ctx *gin.Context) {
	req := entities.NewMatchRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if req.UserID == req.UserMatch {
		res := entities.UnprocessableEntityResponse()
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	userExist := controller.credential.New()
	userExist.ID = req.UserID
	if exist := controller.credential.FindOne(userExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	otherUserExist := controller.credential.New()
	otherUserExist.ID = req.UserMatch
	if exist := controller.credential.FindOne(otherUserExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	matchExist := controller.match.New()
	matchExist.UserID = req.UserID
	matchExist.UserMatch = req.UserMatch
	if duplicate := controller.match.FindOne(matchExist); duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	data := controller.match.New()
	data.UserID = req.UserID
	data.UserMatch = req.UserMatch
	data.Liked = req.Liked

	match, err := controller.match.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewMatchResponse(match.ID, match.UserID, match.UserMatch, match.Liked)
	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *MatchController) Put(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		match := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, match)
		return
	}

	checkExisting := controller.match.New()
	checkExisting.ID = id

	if exist := controller.match.FindOne(checkExisting); exist.ID == 0 {
		match := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, match)
		return
	}

	req := entities.NewMatchRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		match := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, match)
		return
	}

	data := controller.match.New()
	data.ID = id
	data.UserID = req.UserID
	data.UserMatch = req.UserMatch
	data.Liked = req.Liked

	match, err := controller.match.UpdateByID(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewMatchResponse(match.ID, match.UserID, match.UserMatch, match.Liked)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *MatchController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		match := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, match)
		return
	}

	data := controller.match.New()
	data.ID = id

	match := controller.match.FindOne(data)
	if match.ID == 0 {
		match := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, match)
		return
	}

	if err := controller.match.Delete(data); err != nil {
		match := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, match)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
	return
}
