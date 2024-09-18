package sygnal

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

//файл для ловли сигналов завершения, чтобы не потерять данные

func SygnalSearch(oS *storage.OrderStorage, rS *storage.ReturnStorage) error {
	signalls := make(chan os.Signal, 1)

	signal.Notify(signalls, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalls
		fmt.Println()
		fmt.Println("exit")
		oS.WriteToJSON()
		rS.WritoToJSON()
		os.Exit(1)

	}()
	return nil
}
