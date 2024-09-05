package main

import (
	"gitlab.ozon.dev/akugnerevich/homework-1.git/cmd/sygnal"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/basic"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/cli"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"os"
)

func main() {
	strg := storage.NewOrderStorage()
	err := strg.ReadFromJSON("data/orders.json")
	if err != nil {
		return
	}
	sygnal.SygnalSearch(strg)
	if len(os.Args) > 1 {
		cli.Execute()
	} else {
		basic.Run(strg)
	}
}
