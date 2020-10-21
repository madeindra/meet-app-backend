package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/controllers"
	"github.com/madeindra/meet-app/db"
	"github.com/madeindra/meet-app/models"
)

const (
	rootPath string = "/"
	v1Path   string = "/api/v1"

	registerPath string = "/registration"
)

func Init() *gin.Engine {
	router := gin.Default()

	pingController := controllers.NewPingController()
	credentialController := controllers.NewCredentialController(models.NewCredentialImplementation(db.DB))

	router.GET(rootPath, pingController.Ping)

	v1 := router.Group(v1Path)
	{
		v1.POST(registerPath, credentialController.Create)
	}

	return router
}
