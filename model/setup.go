package model

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func Init() {
	conn, err := gorm.Open("sqlite3", ":memory:")

	if err != nil {
		panic("Failed while connecting to database")
	}

	conn.AutoMigrate(&User{})

	DB = conn
}

func Close() {
	DB.Close()
}
