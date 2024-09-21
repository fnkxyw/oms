package Errs

import (
	"errors"
)

var (
	ErrIsConsist       = errors.New("this item is already available at the PuP")
	ErrDate            = errors.New("Incorrect date")
	ErrNotAllIDs       = errors.New("Not all orderStorage are for the same user")
	ErrCheckOrderID    = errors.New("Check input OrderId")
	ErrTimeExpired     = errors.New("Return time has expired :(")
	ErrNotPlace        = errors.New("Order are not palced")
	ErrIncorrectUserId = errors.New("Check input UserID")
)
