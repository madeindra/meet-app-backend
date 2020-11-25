package models

import (
	"github.com/jinzhu/gorm"
)

type resets struct {
	ID           uint64 `gorm:"primaryKey"`
	UserID       uint64 `gorm:"unique"`
	RefreshToken string
}

type ResetInterface interface {
	New() resets
	Create(data resets) (resets, error)
	FindOne(data resets) resets
	UpdateByUser(data resets) (resets, error)
	Delete(data resets) error
}

type ResetImplementation struct {
	db *gorm.DB
}

func NewResetModel(db *gorm.DB) *ResetImplementation {
	return &ResetImplementation{db}
}
