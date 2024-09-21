package sygnal

import (
	"fmt"
	s "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	r "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
	"os"
	"os/signal"
	"syscall"
)

//файл для ловли сигналов завершения, чтобы не потерять данные

func SygnalSearch(oS s.OrderStorageInterface, rS r.ReturnStorageInterface) error {
	signalls := make(chan os.Signal, 1)

	signal.Notify(signalls, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalls
		fmt.Println()
		fmt.Println("exit")
		oS.WriteToJSON()
		rS.WriteToJSON()
		os.Exit(1)

	}()
	return nil
}
