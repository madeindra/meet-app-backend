package models

import "github.com/madeindra/meet-app/common"

func Migrate() {
	common.DB.AutoMigrate(&Credentials{})
}
