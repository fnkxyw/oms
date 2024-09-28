package main

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"math/rand"
)

var states = []models.State{
	models.AcceptState,
	models.RefundedState,
	models.NewState,
	models.PlaceState,
	models.SoftDelete,
}

func generateState() models.State {
	return states[rand.Intn(len(states))]
}
