package controller

import (
	"errors"
	"fmt"
)

func AddWrap() (bool, error) {
	var answer string
	fmt.Println("whether you want to add a wrap to your package?['y' - yes, 'n' - no ]")
	fmt.Scan(&answer)
	switch answer {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, errors.New("Invalid input")

	}

}
