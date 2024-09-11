package errors

import "errors"

var (
	ErrorIsConsist       = errors.New("this item is already available at the PuP")
	ErrorDate            = errors.New("Incorrect date")
	ErrorNotAllIDs       = errors.New("Not all orders are for the same user")
	ErrorCheckOrderID    = errors.New("Check input OrderId")
	ErrorTimeExpired     = errors.New("Return time has expired :(")
	ErrorNotPlace        = errors.New("Order are not palced")
	ErrorIncorrectUserId = errors.New("Check input UserID")
)
