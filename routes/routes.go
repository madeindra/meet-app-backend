package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/controllers"
	"github.com/madeindra/meet-app/helpers"
	"github.com/madeindra/meet-app/middlewares"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/validators"
)

const (
	rootPath string = "/"
	v1Path   string = "/api/v1"

	authenticatePath string = "/authentication"
	registerPath     string = "/registration"
	tokenPath        string = "/token"
	loginPath        string = "/login"
	profilePath      string = "/profiles"
	profileIDPath    string = "/profiles/:id"
)

func RouterInit() *gin.Engine {
	router := gin.Default()
	db := common.GetDB()

	binding.Validator = validators.NewValidator()

	hashHelper := helpers.NewHashHelper()
	bearerHelper := helpers.NewJWTHelper()

	credentialModel := models.NewCredentialModel(db)
	tokenModel := models.NewTokenModel(db)
	profileModel := models.NewProfileModel(db)

	pingController := controllers.NewPingController()
	credentialController := controllers.NewCredentialController(credentialModel, tokenModel, hashHelper, bearerHelper)
	tokenController := controllers.NewTokenController(tokenModel, credentialModel, bearerHelper)
	profileController := controllers.NewProfileController(profileModel, credentialModel)

	router.GET(rootPath, pingController.Ping)

	v1 := router.Group(v1Path)
	auth := v1.Group(authenticatePath)

	auth.Use(middlewares.Basic())

	auth.POST(registerPath, credentialController.Register)
	auth.POST(loginPath, credentialController.Login)
	auth.POST(tokenPath, tokenController.Refresh)

	v1.Use(middlewares.Jwt())

	v1.GET(profileIDPath, profileController.GetSingle)
	v1.PUT(profileIDPath, profileController.Put)
	v1.DELETE(profileIDPath, profileController.Delete)
	v1.GET(profilePath, profileController.GetCollections)
	v1.POST(profilePath, profileController.Post)
	return router
}
