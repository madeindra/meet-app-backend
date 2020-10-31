package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/madeindra/meet-app/common"
)

func Basic() gin.HandlerFunc {
	basicUsername := common.GetBasicUsername()
	basicPassword := common.GetBasicPassword()
	return gin.BasicAuth(gin.Accounts{
		basicUsername: basicPassword,
	})
}
