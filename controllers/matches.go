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
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find existing match in db
	data := controller.match.New()
	data.ID = id

	match := controller.match.FindOne(data)
	if match.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return response
	res := entities.NewMatchResponse(match.ID, match.UserID, match.UserMatch, match.Liked)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *MatchController) GetCollections(ctx *gin.Context) {
	// get userId of the person from url
	userID, err := strconv.ParseUint(ctx.DefaultQuery("userId", "0"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// get the userId of the matched from url
	userMatch, err := strconv.ParseUint(ctx.DefaultQuery("matchTo", "0"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// define variable for boolean sensitive search, without this, it can't filter matched status "false"
	var boolSensitive bool = false
	var liked bool = false

	// if liked query is set, the search is a boolean sensitive search
	if ctx.Query("liked") != "" {
		liked, err = strconv.ParseBool(ctx.Query("liked"))
		if err != nil {
			res := entities.BadRequestResponse()
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		boolSensitive = true
	}

	// create match data to be found in db
	data := controller.match.New()
	data.UserID = userID
	data.UserMatch = userMatch
	data.Liked = liked

	// find match data in db
	match := controller.match.FindBy(data, boolSensitive)
	if len(match) == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// return resposne
	res := entities.NewMatchBatchResponse(match)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *MatchController) Post(ctx *gin.Context) {
	// bind request
	req := entities.NewMatchRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// disable user who want to match with theirself
	if req.UserID == req.UserMatch {
		res := entities.UnprocessableEntityResponse()
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	// find user in db
	userExist := controller.credential.New()
	userExist.ID = req.UserID
	if exist := controller.credential.FindOne(userExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// find matched user in db
	otherUserExist := controller.credential.New()
	otherUserExist.ID = req.UserMatch
	if exist := controller.credential.FindOne(otherUserExist); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// find existing match data
	matchExist := controller.match.New()
	matchExist.UserID = req.UserID
	matchExist.UserMatch = req.UserMatch
	if duplicate := controller.match.FindOne(matchExist); duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// create match data to be inserted to db
	data := controller.match.New()
	data.UserID = req.UserID
	data.UserMatch = req.UserMatch
	data.Liked = req.Liked

	// insert match data to db
	match, err := controller.match.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewMatchResponse(match.ID, match.UserID, match.UserMatch, match.Liked)
	ctx.JSON(http.StatusCreated, res)
	return
}

func (controller *MatchController) Put(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find match by id in db
	checkExisting := controller.match.New()
	checkExisting.ID = id

	if exist := controller.match.FindOne(checkExisting); exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// bind request
	req := entities.NewMatchRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// create match data to be updated in db
	data := controller.match.New()
	data.ID = id
	data.UserID = req.UserID
	data.UserMatch = req.UserMatch
	data.Liked = req.Liked

	// update match data in db
	match, err := controller.match.UpdateByID(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewMatchResponse(match.ID, match.UserID, match.UserMatch, match.Liked)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *MatchController) Delete(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find match by id
	data := controller.match.New()
	data.ID = id

	match := controller.match.FindOne(data)
	if match.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// delete match data
	if err := controller.match.Delete(data); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	ctx.JSON(http.StatusNoContent, nil)
	return
}
