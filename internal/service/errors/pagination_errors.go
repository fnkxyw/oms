package Errs

import (
	"errors"
)

var (
	ErrLimitPage   = errors.New("page and limit must be greater than 0")
	ErrNoMoreItems = errors.New("No more items avilable")
)
