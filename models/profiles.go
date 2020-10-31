package models

import "github.com/jinzhu/gorm"

type profiles struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	UserID    uint64 `json:"userId"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type ProfilesInterface interface {
	Create(data profiles) (profiles, error)
	FindOne(data profiles) profiles
}

type ProfilesImplementation struct {
	db *gorm.DB
}

func NewProfileData(userId uint, firstName string, lastName string) profiles {
	return profiles{UserID: userId, FirstName: firstName, LastName: lastName}
}

func NewProfileImplementation(db *gorm.DB) *ProfilesImplementation {
	return &ProfilesImplementation{db}
}

func (implementation *ProfilesImplementation) Create(data profiles) (profiles, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return profiles{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ProfilesImplementation) FindOne(data profiles) profiles {
	tx := implementation.db
	res := profiles{}

	tx.Where(data).First(&res)

	return res
}
