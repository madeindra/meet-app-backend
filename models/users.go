package models

import "github.com/madeindra/meet-app/db"

type Users struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

func CreateUser(data Users) (Users, error) {
	tx := db.DB.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Users{}, err
	}

	return data, tx.Commit().Error
}
