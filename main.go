package main

import (
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/routes"
)

func main() {
	common.ConfigInit()
	common.DBInit()

	models.Migrate()

	server := routes.RouterInit()
	server.Run(":" + common.GetServerPort())
}
