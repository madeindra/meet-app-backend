package models

import (
	"gorm.io/gorm"
)

type tickets struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"unique"`
	Ticket string
}

type TicketInterface interface {
	New() tickets
	Create(data tickets) (tickets, error)
	FindOne(data tickets) tickets
	UpdateByUser(data tickets) (tickets, error)
}

type TicketImplementation struct {
	db *gorm.DB
}

func NewTicketModel(db *gorm.DB) *TicketImplementation {
	return &TicketImplementation{db}
}

func (implementation *TicketImplementation) New() tickets {
	return tickets{}
}

func (implementation *TicketImplementation) Create(data tickets) (tickets, error) {
	tx := implementation.db.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return tickets{}, err
	}

	return data, tx.Commit().Error
}

func (implementation *TicketImplementation) FindOne(data tickets) tickets {
	res := tickets{}

	implementation.db.Where(data).First(&res)

	return res
}

func (implementation *TicketImplementation) UpdateByUser(data tickets) (tickets, error) {
	tx := implementation.db.Begin()
	res := tickets{UserID: data.UserID}

	if err := tx.Model(tickets{}).Where(res).Updates(&data).Error; err != nil {
		tx.Rollback()
		return tickets{}, err
	}

	return data, tx.Commit().Error
}
