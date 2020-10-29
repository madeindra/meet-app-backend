package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/controllers"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/validators"
)

const (
	rootPath string = "/"
	v1Path   string = "/api/v1"

	registerPath string = "/registration"
)

func RouterInit() *gin.Engine {
	router := gin.Default()

	binding.Validator = validators.NewValidator()

	pingController := controllers.NewPingController()
	credentialController := controllers.NewCredentialController(models.NewCredentialImplementation(common.DB))

	router.GET(rootPath, pingController.Ping)

	v1 := router.Group(v1Path)
	{
		v1.POST(registerPath, credentialController.Create)
	}

	return router
}
