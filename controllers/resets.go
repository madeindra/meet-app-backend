package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
)

type ResetController struct {
	reset      models.ResetInterface
	credential models.CredentialInterface
	hash       helpers.HashInterface
	random     helpers.RandomInterface
}

const tokenLength = 64

func NewResetController(reset models.ResetInterface, credential models.CredentialInterface, hash helpers.HashInterface, random helpers.RandomInterface) *ResetController {
	return &ResetController{reset, credential, hash, random}
}

func (controller *ResetController) Start(ctx *gin.Context) {
	// bind Request
	req := entities.NewResetStartRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// check credential by email in db
	credential := controller.credential.New()
	credential.Email = req.Email

	exist := controller.credential.FindOne(credential)
	if exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// check reset request on db
	data := controller.reset.New()
	data.UserID = exist.ID

	// create a new by user id if not exist
	if exist := controller.reset.FindOne(data); exist.ID == 0 {
		if _, err := controller.reset.Create(data); err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	// update token in db
	data.Token = controller.random.RandomString(tokenLength)

	if _, err := controller.reset.UpdateByUser(data); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewResetStartResponse(data.UserID, data.Token)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ResetController) Complete(ctx *gin.Context) {
	// get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	// get token from url
	token := ctx.Query("token")
	if token == "" {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// bind request
	data := entities.NewResetCompleteRequest()
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find credential data by id in db
	credentialData := controller.credential.New()
	credentialData.ID = id
	credential := controller.credential.FindOne(credentialData)
	if credential.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// check user id & token match the one stored in db
	resetData := controller.reset.New()
	resetData.UserID = id
	resetData.Token = token
	reset := controller.reset.FindOne(resetData)
	if reset.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// generate hash
	hash, err := controller.hash.Generate(data.Password)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// update password
	credential.Password = hash
	if _, err := controller.credential.UpdateByID(credential); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// delete reset request in table
	controller.reset.Delete(resetData)

	// return response
	res := entities.NewResetCompleteResponse(credential.ID)
	ctx.JSON(http.StatusOK, res)
	return
}
