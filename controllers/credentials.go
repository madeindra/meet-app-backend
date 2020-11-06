package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
)

type CredentialController struct {
	credential models.CredentialInterface
	token      models.TokenInterface
	hash       helpers.HashInterface
	bearer     helpers.JWTInterface
}

func NewCredentialController(credential models.CredentialInterface, token models.TokenInterface, hash helpers.HashInterface, bearer helpers.JWTInterface) *CredentialController {
	return &CredentialController{credential, token, hash, bearer}
}

func (controller *CredentialController) Register(ctx *gin.Context) {
	req := entities.NewCredentialRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user := controller.credential.New()
	user.Email = req.Email
	duplicate := controller.credential.FindOne(user)
	if duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	hash, err := controller.hash.Generate(req.Password)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	data := controller.credential.New()
	data.Email = req.Email
	data.Password = hash

	credential, err := controller.credential.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewCredentialResponse(credential.ID, credential.Email)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *CredentialController) Login(ctx *gin.Context) {
	req := entities.NewCredentialRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user := controller.credential.New()
	user.Email = req.Email
	credential := controller.credential.FindOne(user)
	if credential.ID == 0 {
		res := entities.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	err := controller.hash.Verify(credential.Password, req.Password)
	if err != nil {
		res := entities.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	token, err := controller.bearer.GenerateToken(credential.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshToken, err := controller.bearer.GenerateRefresh(credential.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	refreshTokenData := controller.token.New()
	refreshTokenData.UserID = credential.ID
	refreshTokenData.RefreshToken = refreshToken
	if _, err := controller.token.Update(refreshTokenData); err != nil {
		if _, err := controller.token.Create(refreshTokenData); err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	res := entities.NewAuthenticatedResponse(credential.ID, credential.Email, token, refreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}
