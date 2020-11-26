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
	// bind request
	req := entities.NewCredentialRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find existing credential by email
	user := controller.credential.New()
	user.Email = req.Email
	if duplicate := controller.credential.FindOne(user); duplicate.ID != 0 {
		res := entities.ConflictResponse()
		ctx.JSON(http.StatusConflict, res)
		return
	}

	// create hash from password in request
	hash, err := controller.hash.Generate(req.Password)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// create credential data to insert
	data := controller.credential.New()
	data.Email = req.Email
	data.Password = hash

	// insert credential data to db
	credential, err := controller.credential.Create(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// create profile data for that credential
	userProfile := controller.profile.New()
	userProfile.UserID = credential.ID

	// insert profile data to db
	_, err = controller.profile.Create(userProfile)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewCredentialResponse(credential.ID, credential.Email)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *CredentialController) Login(ctx *gin.Context) {
	// bind request
	req := entities.NewCredentialRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find existing credential by email
	user := controller.credential.New()
	user.Email = req.Email
	credential := controller.credential.FindOne(user)
	if credential.ID == 0 {
		res := entities.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	// verify hashed password in db to plain password in request
	err := controller.hash.Verify(credential.Password, req.Password)
	if err != nil {
		res := entities.UnauthorizedResponse()
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	// generate jwt bearer token
	token, err := controller.bearer.GenerateToken(credential.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// generate jwt refresh token
	refreshToken, err := controller.bearer.GenerateRefresh(credential.Email)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// create refresh token data
	refreshTokenData := controller.token.New()
	refreshTokenData.UserID = credential.ID

	// find existing refresh token data in db, if not exist create with user id
	if exist := controller.token.FindOne(refreshTokenData); exist.ID == 0 {
		if _, err := controller.token.Create(refreshTokenData); err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	// assign refresh token to update
	refreshTokenData.RefreshToken = refreshToken

	// update refresh token in db
	if _, err := controller.token.UpdateByUser(refreshTokenData); err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewAuthenticatedResponse(credential.ID, credential.Email, token, refreshToken)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *CredentialController) Update(ctx *gin.Context) {
	// get ID from param
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

	// bind request
	req := entities.NewCredentialUpdateRequest()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := entities.BadRequestResponse()
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// find credential by id and email in database
	data := controller.credential.New()
	data.ID = id
	data.Email = fmt.Sprintf("%v", email)

	credential := controller.credential.FindOne(data)
	if credential.Password == "" {
		res := entities.NotFoundResponse()
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	// if user request password update, check if the password in db match with password in request
	if req.OldPassword != "" && req.NewPassword != "" {
		if err := controller.hash.Verify(credential.Password, req.OldPassword); err != nil {
			res := entities.UnauthorizedResponse()
			ctx.JSON(http.StatusUnauthorized, res)
			return
		}

		// create hash of new password
		hash, err := controller.hash.Generate(req.NewPassword)
		if err != nil {
			res := entities.InterenalServerErrorResponse()
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}

		// assign hash to password
		data.Password = hash
	}

	// if user request email update, assign it
	if req.Email != "" {
		data.Email = req.Email
	}

	// update user credential
	_, err = controller.credential.UpdateByID(data)
	if err != nil {
		res := entities.InterenalServerErrorResponse()
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// return response
	res := entities.NewCredentialUpdateResponse(data.ID, data.Email)
	ctx.JSON(http.StatusOK, res)
	return
}
