package main

import (
	"fmt"
	signals "gitlab.ozon.dev/akugnerevich/homework.git/cmd/signals"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	s "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	r "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
)

func main() {
	orderStorage := s.NewOrderStorage()
	returnStorage := r.NewReturnStorage()
	err := orderStorage.ReadFromJSON()
	if err != nil {
		fmt.Println(err)
		err = orderStorage.Create()
		if err != nil {
			fmt.Println(err)
		}
	}
	err = returnStorage.ReadFromJSON()
	if err != nil {
		fmt.Println(err)
		err = returnStorage.Create()
		if err != nil {
			fmt.Println(err)
		}
	}
	signals.SygnalSearch(orderStorage, returnStorage)

	c.Run(orderStorage, returnStorage)
}
