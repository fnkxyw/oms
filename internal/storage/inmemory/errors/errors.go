package errors

import "errors"

var (
	ErrOpenFile    = errors.New("Open File error ")
	ErrReadFile    = errors.New("Read File error ")
	ErrEnocde      = errors.New("Encoding error ")
	ErrAlrAccept   = errors.New("Order already accept ")
	ErrAlrReturn   = errors.New("Order already return ")
	ErrInvalidType = errors.New("Invalid type error ")
)
