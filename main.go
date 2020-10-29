package main

import (
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/routes"
)

func main() {
	common.DBInit()
	models.Migrate()
	defer common.DBClose()

	server := routes.Init()
	server.Run(":8080")
}
