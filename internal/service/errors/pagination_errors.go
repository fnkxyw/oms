package errors

import "errors"

var (
	ErrorLimitPage   = errors.New("page and limit must be greater than 0")
	ErrorNoMoreItems = errors.New("No more items avilable")
)
