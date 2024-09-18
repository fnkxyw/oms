package Errs

import (
	"errors"
)

var (
	ErrNoConsist = errors.New("We dont have order with that id")
)
