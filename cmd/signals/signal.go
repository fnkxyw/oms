package sygnal

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
	"os"
	"os/signal"
	"syscall"
)

//файл для ловли сигналов завершения, чтобы не потерять данные

func SygnalSearch(oS orderStorage.OrderStorageInterface, rS returnStorage.ReturnStorageInterface) error {
	signalls := make(chan os.Signal, 1)

	signal.Notify(signalls, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalls
		fmt.Println()
		fmt.Println("exit")
		err := oS.WriteToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write order storage to JSON: %v\n", err)

			return
		}

		err = rS.WriteToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write return storage to JSON: %v\n", err)
			return
		}
		os.Exit(1)

	}()
	return nil
}
