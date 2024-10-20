package models

import (
	"time"
)

// модель заказа
type Order struct {
	ID            uint      `json:"id" db:"id"`
	UserID        uint      `json:"user_id" db:"user_id" `
	State         State     `json:"state" db:"state" `
	AcceptTime    int64     `json:"accept_time" db:"accept_time" `
	KeepUntilDate time.Time `json:"date" db:"keep_until_date" `
	PlaceDate     time.Time `json:"place_data" db:"place_date" `
	Weight        int       `json:"weight" db:"weight" `
	Price         int       `json:"price" db:"price"`
}

type State string

// состояния заказа
var (
	SoftDelete    = State("soft_delete")
	AcceptState   = State("accept")
	PlaceState    = State("place")
	RefundedState = State("refunded")
	NewState      = State("new_state")
)

type Event string

var (
	AcceptEvent = Event("accept")
	PlaceEvent  = Event("place")
	ReturnEvent = Event("return")
	RefundEvent = Event("refund")
)
