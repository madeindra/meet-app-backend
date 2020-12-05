package common

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func DBInit() {
	host := GetDatabaseHost()

	conn, err := gorm.Open(sqlite.Open(host), &gorm.Config{})

	if err != nil {
		panic("Failed while connecting to database")
	}

	db = conn
}

func GetDB() *gorm.DB {
	return db
}
