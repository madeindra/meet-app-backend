package middlewares

import "github.com/gin-gonic/gin"

//TODO: Get from env/config
const (
	basicUsername = "foo"
	basicPassword = "bar"
)

func Basic() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		basicUsername: basicPassword,
	})
}
