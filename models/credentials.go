package models

import "github.com/madeindra/meet-app/db"

type Credentials struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func CreateCredential(data Credentials) (Credentials, error) {
	tx := db.DB.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Credentials{}, err
	}

	return data, tx.Commit().Error
}
