package main

import (
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/basic"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/cli"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		cli.Execute()
	} else {
		basic.Run()
	}
}
