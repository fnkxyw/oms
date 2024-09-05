package sygnal

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

func SygnalSearch(s *storage.OrderStorage) error {
	signalls := make(chan os.Signal, 1)

	signal.Notify(signalls, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			<-signalls
			fmt.Println()
			fmt.Println("exit")
			storage.WriteToJSON("data/orders.json", s)
			os.Exit(1)

		}
	}()
	return nil
}
