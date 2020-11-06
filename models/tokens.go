package models

import (
	"github.com/jinzhu/gorm"
)

type tokens struct {
	ID           uint64 `gorm:"primaryKey"`
	UserID       uint64 `gorm:"unique"`
	RefreshToken string
}

type TokenInterface interface {
	New() tokens
	Create(data tokens) (tokens, error)
	FindOne(data tokens) tokens
	UpdateByUser(data tokens) (tokens, error)
}

type TokenImplementation struct {
	db *gorm.DB
}

func NewTokenModel(db *gorm.DB) *TokenImplementation {
	return &TokenImplementation{db}
}

func (implementation *TokenImplementation) New() tokens {
	return tokens{}
}

func (implementation *TokenImplementation) Create(data tokens) (tokens, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return tokens{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *TokenImplementation) FindOne(data tokens) tokens {
	res := tokens{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *TokenImplementation) UpdateByUser(data tokens) (tokens, error) {
	tx := implementation.db.Begin()
	res := tokens{UserID: data.UserID}

	if err := tx.Model(tokens{}).Where(res).Updates(&data).Error; err != nil {
		tx.Rollback()
		return tokens{}, err
	}

	return data, tx.Commit().Error
}
