package models

import "github.com/jinzhu/gorm"

type Matches struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	UserMatch uint64
	Liked     bool
}

type MatchInterface interface {
	New() Matches
	Create(data Matches) (Matches, error)
	FindAll() []Matches
	FindOne(data Matches) Matches
	UpdateByID(data Matches) (Matches, error)
	Delete(data Matches) error
}

type MatchImplementation struct {
	db *gorm.DB
}

func NewMatchModel(db *gorm.DB) *MatchImplementation {
	return &MatchImplementation{db}
}

func (implementation *MatchImplementation) New() Matches {
	return Matches{}
}

func (implementation *MatchImplementation) Create(data Matches) (Matches, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Matches{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *MatchImplementation) FindAll() []Matches {
	res := []Matches{}

	implementation.db.Find(&res)

	return res
}
func (implementation *MatchImplementation) FindOne(data Matches) Matches {
	res := Matches{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *MatchImplementation) UpdateByID(data Matches) (Matches, error) {
	tx := implementation.db.Begin()
	res := Matches{ID: data.ID}

	if err := tx.Model(Matches{}).Where(res).Updates(map[string]interface{}{"liked": data.Liked}).Error; err != nil {
		tx.Rollback()
		return Matches{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *MatchImplementation) Delete(data Matches) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
