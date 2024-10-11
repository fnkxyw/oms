package pup_service

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
)

type Implementation struct {
	storage storage.Facade
	desc.UnimplementedPupServiceServer
}

func NewImplementation(storage storage.Facade) *Implementation {
	return &Implementation{storage: storage}
}
