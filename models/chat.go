package models

import (
	"gorm.io/gorm"
)

type Chats struct {
	ID      uint64 `gorm:"primaryKey" json:",omitempty"`
	Sender  uint64
	Target  uint64
	Content string
}

type ChatsInterface interface {
	New() Chats
	Create(data Chats) (Chats, error)
	FindAll() []Chats
	FindOne(data Chats) Chats
	Delete(data Chats) error
}

type ChatsImplementation struct {
	db *gorm.DB
}

func NewChatModel(db *gorm.DB) *ChatsImplementation {
	return &ChatsImplementation{db}
}

func (implementation *ChatsImplementation) New() Chats {
	return Chats{}
}

func (implementation *ChatsImplementation) Create(data Chats) (Chats, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return Chats{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *ChatsImplementation) FindAll() []Chats {
	res := []Chats{}

	implementation.db.Find(&res)

	return res
}

func (implementation *ChatsImplementation) FindOne(data Chats) Chats {
	res := Chats{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *ChatsImplementation) Delete(data Chats) error {
	tx := implementation.db.Begin()

	if err := tx.Where(data).Delete(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
