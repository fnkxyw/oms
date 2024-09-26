package main

import (
	"fmt"
	signals "gitlab.ozon.dev/akugnerevich/homework.git/cmd/signals"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
)

func main() {
	oS := orderStorage.NewOrderStorage()
	err := oS.ReadFromJSON()
	if err != nil {
		fmt.Println(err)
		err = oS.Create()
		if err != nil {
			fmt.Println(err)
		}
	}
	err = signals.SygnalSearch(*oS)
	if err != nil {
		return
	}

	err = c.Run(oS)
	if err != nil {
		return
	}
}
