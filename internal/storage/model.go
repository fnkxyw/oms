package storage

import "gitlab.ozon.dev/akugnerevich/homework.git/internal/models"

type Storage interface {
	AddToStorage(order *models.Order)
	IsConsist(id uint) bool
	DeleteFromStorage(id uint)
	GetItem(id uint) (*models.Order, bool)
	GetIDs() []uint
}
