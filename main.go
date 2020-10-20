package main

import (
	"github.com/madeindra/meet-app/model"
	"github.com/madeindra/meet-app/route"
)

func main() {
	model.Init()
	defer model.Close()

	server := route.Init()
	server.Run(":8080")
}
