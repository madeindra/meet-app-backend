package main

import (
	"github.com/madeindra/meet-app/db"
	"github.com/madeindra/meet-app/models"
	"github.com/madeindra/meet-app/routes"
)

func main() {
	db.Init()
	models.Migrate()
	defer db.Close()

	server := routes.Init()
	server.Run(":8080")
}
