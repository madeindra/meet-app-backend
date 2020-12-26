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

	chatPath         string = "/chat"
	chatDetailPath   string = "/chat/details"
	chatDetailIDPath string = "/chat/details/:id"
	authenticatePath string = "/authentication"
	registerPath     string = "/registration"
	credentialIDPath string = "/credential/:id"
	tokenPath        string = "/token"
	loginPath        string = "/login"
	resetPath        string = "/reset"
	resetIDPath      string = "/reset/:id"
	profilePath      string = "/profiles"
	profileIDPath    string = "/profiles/:id"
	skillPath        string = "/skills"
	skillIDPath      string = "/skills/:id"
	matchPath        string = "/matches"
	matchIDPath      string = "/matches/:id"
)

func RouterInit() *gin.Engine {
	router := gin.Default()
	db := common.GetDB()

	binding.Validator = validators.NewValidator()

	hashHelper := helpers.NewHashHelper()
	bearerHelper := helpers.NewJWTHelper()
	randomHelper := helpers.NewRandomHelper()

	pubSub := models.NewPubSub()

	pubSubModel := models.NewPubSubModel(pubSub)
	chatModel := models.NewChatModel(db)
	credentialModel := models.NewCredentialModel(db)
	tokenModel := models.NewTokenModel(db)
	ticketModel := models.NewTicketModel(db)
	resetModel := models.NewResetModel(db)
	profileModel := models.NewProfileModel(db)
	skillModel := models.NewSkillModel(db)
	matchModel := models.NewMatchModel(db)

	pingController := controllers.NewPingController()
	pubSubController := controllers.NewPubSubController(pubSubModel, chatModel, ticketModel, credentialModel)
	chatController := controllers.NewChatController(chatModel, credentialModel)
	credentialController := controllers.NewCredentialController(credentialModel, tokenModel, ticketModel, profileModel, hashHelper, bearerHelper, randomHelper)
	resetController := controllers.NewResetController(resetModel, credentialModel, hashHelper, randomHelper)
	tokenController := controllers.NewTokenController(tokenModel, credentialModel, bearerHelper)
	profileController := controllers.NewProfileController(profileModel, credentialModel, skillModel)
	skillController := controllers.NewSkillController(skillModel, profileModel)
	matchController := controllers.NewMatchController(matchModel, credentialModel)

	router.GET(rootPath, pingController.Ping)
	router.GET(chatPath, pubSubController.WebsocketHandler)

	v1 := router.Group(v1Path)
	auth := v1.Group(authenticatePath)

	auth.Use(middlewares.Basic())

	auth.POST(registerPath, credentialController.Register)
	auth.POST(loginPath, credentialController.Login)
	auth.POST(tokenPath, tokenController.Refresh)
	auth.POST(resetPath, resetController.Start)
	auth.PUT(resetIDPath, resetController.Complete)

	v1.Use(middlewares.Jwt())

	v1.GET(chatDetailIDPath, chatController.GetDetail)
	v1.GET(chatDetailPath, chatController.GetLatest)

	v1.PUT(credentialIDPath, credentialController.Update)
	v1.GET(profileIDPath, profileController.GetSingle)
	v1.PUT(profileIDPath, profileController.Put)
	v1.DELETE(profileIDPath, profileController.Delete)
	v1.GET(profilePath, profileController.GetCollections)
	v1.POST(profilePath, profileController.Post)

	v1.GET(skillIDPath, skillController.GetSingle)
	v1.PUT(skillIDPath, skillController.Put)
	v1.DELETE(skillIDPath, skillController.Delete)
	v1.GET(skillPath, skillController.GetCollections)
	v1.POST(skillPath, skillController.Post)

	v1.GET(matchIDPath, matchController.GetSingle)
	v1.PUT(matchIDPath, matchController.Put)
	v1.DELETE(matchIDPath, matchController.Delete)
	v1.GET(matchPath, matchController.GetCollections)
	v1.POST(matchPath, matchController.Post)
	return router
}
