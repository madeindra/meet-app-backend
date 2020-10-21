package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func Init() {
	conn, err := gorm.Open("sqlite3", ":memory:")

	if err != nil {
		panic("Failed while connecting to database")
	}

	DB = conn
}

func Close() {
	DB.Close()
}
