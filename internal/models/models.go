package models

import "time"

type Order struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Accept     bool      `json:"accept"`
	AcceptTime time.Time `json:"accept_time"`
	Date       time.Time `json:"date"`
	Issued     bool      `json:"issued"`
	IssuedDate time.Time `json:"issued_date"`
}

type Return struct {
	OrderID      uint      `json:"order_id"`
	UserID       uint      `json:"user_id"`
	DateOfReturn time.Time `json:"date_of_return"`
}
