package storage

import "hw1/internal/models"

type Storage struct {
	Orders  []models.Order
	Returns []models.Return
}
