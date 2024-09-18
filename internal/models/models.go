package models

import "time"

// модель заказа
type Order struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	State         State     `json:"state"`
	AcceptTime    time.Time `json:"accept_time"`
	KeepUntilDate time.Time `json:"date"`
	PlaceDate     time.Time `json:"place_data"`
	Weight        int       `json:"weight"`
	Price         int       `json:"price"`
}

// модель возврата
type Return struct {
	ID           uint      `json:"order_id"`
	UserID       uint      `json:"user_id"`
	DateOfReturn time.Time `json:"date_of_return"`
}

type State string

// состояния заказа
var (
	SoftDelete    = State("SoftDelete")
	AcceptState   = State("accept")
	PlaceState    = State("place")
	ReturnedState = State("returned")
)
