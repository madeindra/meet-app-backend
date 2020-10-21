package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/controllers"
)

const (
	rootPath string = "/"
	v1Path   string = "/api/v1"

	registerPath string = "/registration"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.GET(rootPath, controllers.Ping)

	v1 := router.Group(v1Path)
	{
		v1.POST(registerPath, controllers.CreateCredential)
	}

	return router
}
