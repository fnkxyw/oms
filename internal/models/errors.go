package models

import "errors"

var (
	ErrorCanReturned = errors.New("Order can`t be returned because it has not expired yet and this order is not a return ")
)
