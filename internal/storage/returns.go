package storage

import "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"

type ReturnStorage struct {
	Data map[uint][]models.Return
}

func NewReturnStorage() *ReturnStorage {
	return &ReturnStorage{Data: make(map[uint][]models.Return)}
}
