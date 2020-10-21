package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/controllers"
)

const (
	rootPath string = "/"
	v1Path   string = "/api/v1"

	userPath string = "/users"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.GET(rootPath, controllers.Ping)

	v1 := router.Group(v1Path)
	{
		v1.POST(userPath, controllers.UserCreate)
	}

	return router
}
