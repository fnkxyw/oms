package controller

import (
	"errors"
	"fmt"
)

//пришел к решению что так необходимо сделать чтобы корректно протестировать код

type WrapAdder interface {
	AddWrap() (bool, error)
}

type defaultWrapAdder struct{}

func (d *defaultWrapAdder) AddWrap() (bool, error) {
	var answer string
	fmt.Println("Would you like to add a wrap to your package? ['y' - yes, 'n' - no]")
	fmt.Scan(&answer)
	switch answer {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, errors.New("invalid input")
	}
}

var wrapAdder WrapAdder = &defaultWrapAdder{}

func GetWrapAdder() WrapAdder {
	return wrapAdder
}

func SetWrapAdder(wa WrapAdder) {
	wrapAdder = wa
}
