package models

import "github.com/jinzhu/gorm"

type credentials struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CredentialInterface interface {
	New() credentials
	Create(data credentials) (credentials, error)
	FindOne(data credentials) credentials
}

type CredentialImplementation struct {
	db *gorm.DB
}

func NewCredentialModel(db *gorm.DB) *CredentialImplementation {
	return &CredentialImplementation{db}
}

func (implementation *CredentialImplementation) New() credentials {
	return credentials{}
}

func (implementation *CredentialImplementation) Create(data credentials) (credentials, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return credentials{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *CredentialImplementation) FindOne(data credentials) credentials {
	res := credentials{}

	implementation.db.Where(data).First(&res)

	return res
}
