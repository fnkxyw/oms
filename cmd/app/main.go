package main

import (
	signals "gitlab.ozon.dev/akugnerevich/homework-1.git/cmd/signals"
	c "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/cli"
	s "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
)

func main() {
	orderStorage := s.NewOrderStorage()
	returnStorage := s.NewReturnStorage()
	err := orderStorage.ReadFromJSON("api/orders.json")
	if err != nil {
		return
	}
	err = returnStorage.ReadFromJSON("api/returns.json")
	if err != nil {
		return
	}
	signals.SygnalSearch(orderStorage, returnStorage)

	c.Run(orderStorage, returnStorage)

}
