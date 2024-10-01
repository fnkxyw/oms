package signals

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
	"os"
	"os/signal"
	"syscall"
)

//файл для ловли сигналов завершения, чтобы не потерять данные

func SignalSearch(oS orderStorage.OrderStorage, done chan struct{}) {
	signals := make(chan os.Signal, 1)

	// Подписка на SIGINT и SIGTERM
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println()
		fmt.Println("Received interrupt signal, exiting...")
		if err := oS.WriteToJSON(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write order storage to JSON: %v\n", err)
		}
		close(done) // сигнализируем о завершении
	}()
}
