package pup_service

import (
	"google.golang.org/grpc/codes"
)

var (
	ErrIsPacked = &CustomError{"Order is already packed", codes.FailedPrecondition}
	ErrNotFound = &CustomError{"Order not found", codes.NotFound}
	ErrInternal = &CustomError{"Internal server error", codes.Internal}
)

type CustomError struct {
	Message string
	Code    codes.Code
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) GRPCStatus() codes.Code {
	return e.Code
}
