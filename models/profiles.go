package models

import (
	"github.com/jinzhu/gorm"
)

type Profiles struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	UserID      uint64  `json:"userId" binding:"required" gorm:"unique"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type ProfilesInterface interface {
	New(userId uint64, firstName string, lastName string, description string, latitude float64, longitude float64) Profiles
	Create(data Profiles) (Profiles, error)
	FindAll() []Profiles
	FindByUser(data Profiles) Profiles
	UpdateByUser(data Profiles) (Profiles, error)
	DeleteByUser(data Profiles) error
}

type ProfilesImplementation struct {
	db *gorm.DB
}

func NewProfileModel(db *gorm.DB) *ProfilesImplementation {
	return &ProfilesImplementation{db}
}

func (implementation *ProfilesImplementation) New(userId uint64, firstName string, lastName string, description string, latitude float64, longitude float64) Profiles {
	return Profiles{UserID: userId, FirstName: firstName, LastName: lastName, Description: description, Latitude: latitude, Longitude: longitude}
}

func (implementation *ProfilesImplementation) Create(data Profiles) (Profiles, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Profiles{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ProfilesImplementation) FindAll() []Profiles {
	res := []Profiles{}

	implementation.db.Find(&res)

	return res
}

func (implementation *ProfilesImplementation) FindByUser(data Profiles) Profiles {
	res := Profiles{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *ProfilesImplementation) UpdateByUser(data Profiles) (Profiles, error) {
	tx := implementation.db.Begin()
	res := Profiles{UserID: data.UserID}

	if err := tx.Model(Profiles{}).Where(res).Updates(&data).Error; err != nil {
		tx.Rollback()
		return Profiles{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ProfilesImplementation) DeleteByUser(data Profiles) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
