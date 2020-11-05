package main

import (
	"github.com/madeindra/meet-app/common"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/routes"
)

func main() {
	common.ConfigInit()
	common.DBInit()
	defer common.DBClose()

	models.Migrate()

	server := routes.RouterInit()
	server.Run(":" + common.GetServerPort())
}

//TODO: Create Request Body & NewRequest() for Credential & Token
//TODO: Change New() Model for Credential & Token
//TODO: Update Response Body for Credential & Token
//TODO: Update Model following profiles (don't include primary key) (change function)
//TODO: Update Controller following profiles (create new request & set values) (bind to request body type)
