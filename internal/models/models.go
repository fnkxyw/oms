package models

import "time"

type Order struct {
	ID         uint      `json:"id"`
	CustomerID uint      `json:"customer_id"`
	Duration   time.Time `json:"expiration"`
	Issued     bool      `json:"issued"`
	IssuedDate time.Time `json:"issued_date"`
}

type Return struct {
	OrderID      uint      `json:"order_id"`
	CustomerID   uint      `json:"customer_id"`
	DateOfReturn time.Time `json:"date_of_return"`
}

func NewOrder(id, customerID uint, duration time.Time) *Order {
	return &Order{
		ID:         id,
		CustomerID: customerID,
		Duration:   duration,
		Issued:     false,
	}
}
