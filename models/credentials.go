package models

import "github.com/jinzhu/gorm"

type credentials struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CredentialInterface interface {
	CreateNewCredential(data credentials) (credentials, error)
}

type CredentialImplementation struct {
	db *gorm.DB
}

func NewCredentialData() credentials {
	return credentials{}
}

func NewCredentialImplementation(db *gorm.DB) *CredentialImplementation {
	return &CredentialImplementation{db}
}

func (implementation *CredentialImplementation) CreateNewCredential(data credentials) (credentials, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return credentials{}, err
	}

	return data, tx.Commit().Error
}
