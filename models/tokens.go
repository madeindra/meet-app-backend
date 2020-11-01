package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type tokens struct {
	ID           uint64 `json:"id" gorm:"primaryKey"`
	UserID       uint64 `json:"userId" gorm:"unique"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenInterface interface {
	New(userID uint64, refreshToken string) tokens
	Create(data tokens) (tokens, error)
	FindOne(data tokens) tokens
	Update(data tokens) (tokens, error)
}

type TokenImplementation struct {
	db *gorm.DB
}

func NewTokenImplementation(db *gorm.DB) *TokenImplementation {
	return &TokenImplementation{db}
}

func (implementation *TokenImplementation) New(userID uint64, refreshToken string) tokens {
	return tokens{UserID: userID, RefreshToken: refreshToken}
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
	tx := implementation.db
	res := tokens{}

	tx.Where(data).First(&res)

	return res
}

func (implementation *TokenImplementation) Update(data tokens) (tokens, error) {
	tx := implementation.db.Begin()
	res := tokens{}

	old := implementation.New(data.ID, "")

	if err := tx.Model(&old).Updates(data).Find(&res).Error; err != nil {
		tx.Rollback()
		return tokens{}, err
	}

	if res.ID == 0 {
		return tokens{}, errors.New("Token not found")
	}

	return res, tx.Commit().Error
}
