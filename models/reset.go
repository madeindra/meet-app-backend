package models

import (
	"github.com/jinzhu/gorm"
)

type resets struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"unique"`
	Token  string
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

func (implementation *ResetImplementation) New() resets {
	return resets{}
}

func (implementation *ResetImplementation) Create(data resets) (resets, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return resets{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ResetImplementation) FindOne(data resets) resets {
	res := resets{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *ResetImplementation) UpdateByUser(data resets) (resets, error) {
	tx := implementation.db.Begin()
	res := resets{UserID: data.UserID}

	if err := tx.Model(resets{}).Where(res).Updates(&data).Error; err != nil {
		tx.Rollback()
		return resets{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ResetImplementation) Delete(data resets) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
