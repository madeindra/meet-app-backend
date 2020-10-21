package model

import "github.com/madeindra/meet-app/db"

func Migrate() {
	db.DB.AutoMigrate(&User{})
}
