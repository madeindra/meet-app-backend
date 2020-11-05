package models

import "github.com/madeindra/meet-app/common"

func Migrate() {
	db := common.GetDB()
	db.AutoMigrate(&credentials{}, &Profiles{}, &tokens{})
}
