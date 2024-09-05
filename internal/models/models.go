package models

import "time"

type Order struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	State      State     `json:"state"`
	AcceptTime time.Time `json:"accept_time"`
	Date       time.Time `json:"date"`
	IssuedDate time.Time `json:"issued_date"`
}

type Return struct {
	ID           uint      `json:"order_id"`
	UserID       uint      `json:"user_id"`
	DateOfReturn time.Time `json:"date_of_return"`
}

type State string

var (
	AcceptState   = State("accept")
	PlaceState    = State("place")
	ReturnedState = State("returned")
)
