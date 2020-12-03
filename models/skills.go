package models

import "github.com/jinzhu/gorm"

type Skills struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64
	Name   string
}

type SkillInterface interface {
	New() Skills
	Create(data Skills) (Skills, error)
	FindAll() []Skills
	FindBy(data Skills) []Skills
	FindOne(data Skills) Skills
	UpdateByUser(data Skills) (Skills, error)
	Delete(data Skills) error
}

type SkillImplementation struct {
	db *gorm.DB
}

func NewSkillModel(db *gorm.DB) *SkillImplementation {
	return &SkillImplementation{db}
}

func (implementation *SkillImplementation) New() Skills {
	return Skills{}
}

func (implementation *SkillImplementation) Create(data Skills) (Skills, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Skills{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *SkillImplementation) FindAll() []Skills {
	res := []Skills{}

	implementation.db.Find(&res)

	return res
}

func (implementation *SkillImplementation) FindBy(data Skills) []Skills {
	res := []Skills{}

	implementation.db.Where(data).Find(&res)
	return res
}

func (implementation *SkillImplementation) FindOne(data Skills) Skills {
	res := Skills{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *SkillImplementation) UpdateByUser(data Skills) (Skills, error) {
	tx := implementation.db.Begin()
	res := Skills{UserID: data.UserID}

	if err := tx.Model(Skills{}).Where(res).Updates(&data).Error; err != nil {
		tx.Rollback()
		return Skills{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *SkillImplementation) Delete(data Skills) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
