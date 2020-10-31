package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func DBInit() {
	conn, err := gorm.Open("sqlite3", ":memory:")

	if err != nil {
		panic("Failed while connecting to database")
	}

	db = conn
}

func DBClose() {
	db.Close()
}

func GetDB() *gorm.DB {
	return db
}
