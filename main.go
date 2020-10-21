package main

import (
	"github.com/madeindra/meet-app/db"
	"github.com/madeindra/meet-app/model"
	"github.com/madeindra/meet-app/route"
)

func main() {
	db.Init()
	model.Migrate()
	defer db.Close()

	server := route.Init()
	server.Run(":8080")
}
