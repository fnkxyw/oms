package packing

import "errors"

var (
	ErrWeightBox    = errors.New("weight of goods exceeds 30 kilograms, you can choose only Wrap packaging")
	ErrWeightBundle = errors.New("weight of goods excedes 10 kilograms, you can choose Box or Wrap packing")
	ErrInvalidType  = errors.New("invalid packaging type")
	ErrIsPackaged   = errors.New("order already packaged, you can add wrap")
)
