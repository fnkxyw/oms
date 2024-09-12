package packing

import "errors"

var (
	ErrorWeightBox    = errors.New("weight of goods exceeds 30 kilograms, you can choose only Wrap packaging")
	ErrorWeightBundle = errors.New("weight of goods excedes 10 kilograms, you can choose Box or Wrap packing")
	ErrorInvalidType  = errors.New("invalid packaging type")
	ErrorIsPackaged   = errors.New("order already packaged, you can add wrap")
)
