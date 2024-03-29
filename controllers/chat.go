package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/models"
)

type ChatController struct {
	chat       models.ChatsInterface
	credential models.CredentialInterface
}

func NewChatController(chat models.ChatsInterface, credential models.CredentialInterface) *ChatController {
	return &ChatController{chat, credential}
}

func (controller *ChatController) GetDetail(ctx *gin.Context) {
	// get target ID from param
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// get email from jwt
	email, set := ctx.Get("sub")
	if !set {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// find user by id and email in database
	data := controller.credential.New()
	data.Email = fmt.Sprintf("%v", email)

	user := controller.credential.FindOne(data)
	if user.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// find chat where above user is the sender
	chat := controller.chat.New()
	chat.SenderID = user.ID

	// also find where id in param is the target
	chat.TargetID = id
	details := controller.chat.FindBy(chat)

	// return without checking data length
	res := entities.NewChatBatchResponse(details)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *ChatController) GetLatest(ctx *gin.Context) {
	// get email from jwt
	email, set := ctx.Get("sub")
	if !set {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// find user by id and email in database
	data := controller.credential.New()
	data.Email = fmt.Sprintf("%v", email)

	user := controller.credential.FindOne(data)
	if user.ID == 0 {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// find latest chat where above user is the sender
	chat := controller.chat.New()
	chat.SenderID = user.ID

	details := controller.chat.FindDistinct(chat)

	// return without checking data length
	res := entities.NewChatBatchResponse(details)
	ctx.JSON(http.StatusOK, res)
	return

}
