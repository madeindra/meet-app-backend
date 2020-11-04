package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type profiles struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	UserID      uint64  `json:"userId" binding:"required"`
	FirstName   string  `json:"firstName" binding:"required"`
	LastName    string  `json:"lastName" binding:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type ProfilesInterface interface {
	New(userId uint64, firstName string, lastName string, description string, latitude float64, longitude float64) profiles
	Create(data profiles) (profiles, error)
	FindAll() []profiles
	FindByUser(data profiles) profiles
	UpdateByUser(data profiles) (profiles, error)
	DeleteByUser(data profiles) error
}

type ProfilesImplementation struct {
	db *gorm.DB
}

func NewProfileModel(db *gorm.DB) *ProfilesImplementation {
	return &ProfilesImplementation{db}
}

func (implementation *ProfilesImplementation) New(userId uint64, firstName string, lastName string, description string, latitude float64, longitude float64) profiles {
	return profiles{UserID: userId, FirstName: firstName, LastName: lastName, Description: description, Latitude: latitude, Longitude: longitude}
}

func (implementation *ProfilesImplementation) Create(data profiles) (profiles, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return profiles{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ProfilesImplementation) FindAll() []profiles {
	res := []profiles{}

	implementation.db.Find(&res)

	return res
}

func (implementation *ProfilesImplementation) FindByUser(data profiles) profiles {
	res := profiles{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *ProfilesImplementation) UpdateByUser(data profiles) (profiles, error) {
	tx := implementation.db.Begin()
	res := profiles{}

	if err := tx.First(&res).Updates(data).Error; err != nil {
		tx.Rollback()
		return profiles{}, err
	}

	if res.ID == 0 {
		return profiles{}, errors.New("Not found")
	}

	return res, tx.Commit().Error
}

func (implementation *ProfilesImplementation) DeleteByUser(data profiles) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
