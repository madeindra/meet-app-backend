package models

import (
	"github.com/jinzhu/gorm"
)

type Profiles struct {
	ID          uint64 `gorm:"primaryKey"`
	UserID      uint64 `gorm:"unique"`
	FirstName   string
	LastName    string
	Description string
	Gender      string
	Latitude    float64
	Longitude   float64
}

type ProfilesInterface interface {
	New() Profiles
	Create(data Profiles) (Profiles, error)
	FindAll() []Profiles
	FindOne(data Profiles) Profiles
	UpdateByUser(data Profiles) (Profiles, error)
	Delete(data Profiles) error
}

type ProfilesImplementation struct {
	db *gorm.DB
}

func NewProfileModel(db *gorm.DB) *ProfilesImplementation {
	return &ProfilesImplementation{db}
}

func (implementation *ProfilesImplementation) New() Profiles {
	return Profiles{}
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

func (implementation *ProfilesImplementation) FindOne(data Profiles) Profiles {
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

func (implementation *ProfilesImplementation) Delete(data Profiles) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
