package main

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"math/rand"
)

func generateState() models.State {
	statements := []models.State{
		models.AcceptState,
		models.RefundedState,
		models.NewState,
		models.PlaceState,
		models.SoftDelete,
	}
	return statements[rand.Intn(len(statements))]
}
