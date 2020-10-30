package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/controllers"
	"github.com/madeindra/meet-app/middlewares"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/validators"
)

const (
	rootPath string = "/"
	v1Path   string = "/api/v1"

	registerPath     string = "/registration"
	authenticatePath string = "/authentication"
	profilePath      string = "/profiles"
	profileIdPath    string = "/profiles/:id"
)

func RouterInit() *gin.Engine {
	router := gin.Default()

	binding.Validator = validators.NewValidator()

	pingController := controllers.NewPingController()
	credentialController := controllers.NewCredentialController(models.NewCredentialImplementation(common.DB))
	profileController := controllers.NewProfileController(models.NewProfileImplementation(common.DB))

	router.GET(rootPath, pingController.Ping)

	v1 := router.Group(v1Path)
	{
		v1.POST(registerPath, credentialController.Register)
		v1.POST(authenticatePath, credentialController.Login)

		v1.Use(middlewares.Jwt())

		v1.GET(profilePath, profileController.GetCollections)
		v1.GET(profileIdPath, profileController.GetSingle)
		v1.POST(profilePath, profileController.Post)
		v1.PUT(profilePath, profileController.Put)
		v1.DELETE(profilePath, profileController.Delete)
	}

	return router
}
