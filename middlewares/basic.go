package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/common"
)

func Basic() gin.HandlerFunc {
	var (
		basicUsername string = common.GetBasicUsername()
		basicPassword string = common.GetBasicPassword()
	)

	return gin.BasicAuth(gin.Accounts{
		basicUsername: basicPassword,
	})
}
