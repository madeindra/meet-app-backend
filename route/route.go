package route

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/handler"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.GET("/", handler.Ping)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/user", handler.UserCreate)
	}

	return router
}
