package model

import "github.com/madeindra/meet-app/db"

type User struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

func CreateUser(data User) (User, error) {
	tx := db.DB.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return User{}, err
	}

	return data, tx.Commit().Error
}
