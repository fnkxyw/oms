package main

import (
	"fmt"
	signals "gitlab.ozon.dev/akugnerevich/homework-1.git/cmd/signals"
	c "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/cli"
	s "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
)

func main() {
	orderStorage := s.NewOrderStorage()
	returnStorage := s.NewReturnStorage()
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
