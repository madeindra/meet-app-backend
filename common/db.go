package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func DBInit() {
	provider := GetDatabaseProvider()
	host := GetDatabaseHost()

	conn, err := gorm.Open(provider, host)

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
