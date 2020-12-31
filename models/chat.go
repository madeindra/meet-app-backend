package models

import (
	"gorm.io/gorm"
)

type Chats struct {
	ID            uint64 `gorm:"primaryKey" json:",omitempty"`
	SenderID      uint64
	TargetID      uint64
	Content       string
	SenderName    string
	SenderPicture string
}

type ChatsInterface interface {
	New() Chats
	Create(data Chats) (Chats, error)
	FindBy(data Chats) []Chats
	FindDistinct(data Chats) []Chats
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

func (implementation *ChatsImplementation) FindBy(data Chats) []Chats {
	res := []Chats{}

	implementation.db.Raw("SELECT * FROM chats WHERE sender_id = ? AND target_id = ? UNION SELECT * FROM chats WHERE target_id = ? AND sender_id = ?", data.SenderID, data.TargetID, data.SenderID, data.TargetID).Scan(&res)
	return res
}

func (implementation *ChatsImplementation) FindDistinct(data Chats) []Chats {
	res := []Chats{}

	implementation.db.Raw("SELECT t3.*, t4.name as sender_name, t4.picture as sender_picture FROM(SELECT t1.* FROM chats AS t1 INNER JOIN (SELECT MIN(sender_id, target_id) AS sender_id, MAX(sender_id, target_id) AS target_id, MAX(id) AS max_id FROM chats GROUP BY MIN(sender_id, target_id), MAX(sender_id, target_id)) AS t2 ON MIN(t1.sender_id, t1.target_id) = t2.sender_id AND MAX(t1.sender_id, t1.target_id) = t2.target_id AND t1.id = t2.max_id WHERE t1.sender_id = ? OR t1.target_id = ?) as t3 INNER JOIN profiles as t4 WHERE t3.sender_id = t4.id", data.SenderID, data.SenderID).Scan(&res)

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
