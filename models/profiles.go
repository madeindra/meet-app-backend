package models

import (
	"gorm.io/gorm"
)

type Profiles struct {
	ID           uint64 `gorm:"primaryKey"`
	CredentialID uint64 `gorm:"unique"`
	Name         string
	Description  string
	Gender       string
	Latitude     float64
	Longitude    float64
	Credential   credentials `gorm:"foreignKey:CredentialID"`
	Skills       []Skills    `gorm:"foreignKey:UserID"`
	Matches      []Matches   `gorm:"foreignKey:UserID"`
	Ticket       tickets     `gorm:"foreignKey:UserID"`
	Token        tokens      `gorm:"foreignKey:UserID"`
}

type ProfilesInterface interface {
	New() Profiles
	Create(data Profiles) (Profiles, error)
	FindAll() []Profiles
	FindOne(data Profiles) Profiles
	UpdateByID(data Profiles) (Profiles, error)
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

func (implementation *ProfilesImplementation) UpdateByID(data Profiles) (Profiles, error) {
	tx := implementation.db.Begin()
	res := Profiles{ID: data.ID}

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
