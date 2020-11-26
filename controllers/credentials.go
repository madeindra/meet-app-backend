package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/entities"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/models"
)

type CredentialController struct {
	credential models.CredentialInterface
	token      models.TokenInterface
	profile    models.ProfilesInterface
	hash       helpers.HashInterface
	bearer     helpers.JWTInterface
}

func NewCredentialController(credential models.CredentialInterface, token models.TokenInterface, profile models.ProfilesInterface, hash helpers.HashInterface, bearer helpers.JWTInterface) *CredentialController {
	return &CredentialController{credential, token, profile, hash, bearer}
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
	if duplicate := controller.credential.FindOne(user); duplicate.ID != 0 {
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

	userProfile := controller.profile.New()
	userProfile.UserID = credential.ID

	_, err = controller.profile.Create(userProfile)
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

	if exist := controller.token.FindOne(refreshTokenData); exist.ID == 0 {
		if _, err := controller.token.Create(refreshTokenData); err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	refreshTokenData.RefreshToken = refreshToken

	if _, err := controller.token.UpdateByUser(refreshTokenData); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewAuthenticatedResponse(credential.ID, credential.Email, token, refreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *CredentialController) Update(ctx *gin.Context) {
	// Get ID from param
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get Email from jwt
	email, set := ctx.Get("sub")
	if !set {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Bind Data
	data := entities.NewCredentialUpdateRequest()
	if err := ctx.ShouldBindJSON(&data); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Find ID & Email in database
	credential := controller.credential.New()
	credential.ID = id
	credential.Email = fmt.Sprintf("%v", email)

	exist := controller.credential.FindOne(credential)
	if exist.Password == "" {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// Only if user want to update password
	if data.OldPassword != "" && data.NewPassword != "" {
		// Match old password
		if err := controller.hash.Verify(exist.Password, data.OldPassword); err != nil {
			res := entities.UnauthorizedResponse()
			ctx.JSON(http.StatusUnauthorized, res)
			return
		}

		// create hash of new pass
		hash, err := controller.hash.Generate(data.NewPassword)
		if err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}

		// Update password
		credential.Password = hash
	}

	// Only if user want to update email
	if data.Email != "" {
		credential.Email = data.Email
	}

	_, err = controller.credential.UpdateByID(credential)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := entities.NewCredentialUpdateResponse(credential.ID, credential.Email)
	ctx.JSON(http.StatusOK, res)
	return
}
