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
	// Bind Request
	req := entities.NewResetStartRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Check User Exist
	credential := controller.credential.New()
	credential.Email = req.Email

	exist := controller.credential.FindOne(credential)
	if exist.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Check reset request on table exist / not
	data := controller.reset.New()
	data.UserID = exist.ID

	// insert if not exist
	if exist := controller.reset.FindOne(data); exist.ID == 0 {
		if _, err := controller.reset.Create(data); err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	// update token in table
	data.Token = controller.random.RandomString(tokenLength)

	if _, err := controller.reset.UpdateByUser(data); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewResetStartResponse(data.UserID, data.Token)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ResetController) Complete(ctx *gin.Context) {
	// Get id from url
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		profile := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, profile)
		return
	}

	// Get token from url
	token := ctx.Query("token")
	if token == "" {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Bind Request
	data := entities.NewResetCompleteRequest()
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find email
	credentialData := controller.credential.New()
	credentialData.ID = id
	credential := controller.credential.FindOne(credentialData)
	if credential.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Check user & token on table match / not
	resetData := controller.reset.New()
	resetData.UserID = id
	resetData.Token = token
	reset := controller.reset.FindOne(resetData)
	if reset.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	//create hash
	hash, err := controller.hash.Generate(data.Password)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Update password
	credential.Password = hash
	if _, err := controller.credential.UpdateByID(credential); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Delete / invalidate reset request in table
	controller.reset.Delete(resetData)

	res := entities.NewResetCompleteResponse(credential.ID)
	ctx.JSON(http.StatusOK, res)
	return
}
