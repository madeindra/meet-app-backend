package models

import "github.com/jinzhu/gorm"

type Credentials struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type CredentialInterface interface {
	CreateNewCredential(data Credentials) (Credentials, error)
}

type CredentialImplementation struct {
	db *gorm.DB
}

func NewCredentialImplementation(db *gorm.DB) *CredentialImplementation {
	return &CredentialImplementation{db}
}

func (implementation *CredentialImplementation) CreateNewCredential(data Credentials) (Credentials, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Credentials{}, err
	}

	return data, tx.Commit().Error
}
