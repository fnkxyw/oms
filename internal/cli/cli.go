package cli

import (
	"bufio"
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/wpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var helpText = `
    here is the available set of commands
    help - shows the available commands
    acceptOrder - allows you to take the order from the courier 
    returnOrder - allows you to return the order to the courier
    placeOrder - allow the order to be released to the user
    listOrders - allows you to get a list of orders
    refundOrder - allows you to accept a return from a user
    listReturns - allows you to get a list of returns
    workers-num - allows you to change the number of workers
`

const num_workers = 7

func Run(ctx context.Context, oS storage.Facade) error {
	showHelp()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	notification := make(chan string, 1)

	wg.Add(1)
	go listenChannels(ctx, errChan, notification, &wg)

	wp, err := wpool.NewWorkerPool(ctx, num_workers, notification)
	if err != nil {
		log.Fatalln(err)
	}

	wp.Start()

	wg.Add(1)
	go signalSearch(ctx, cancel, &wg)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		select {
		case <-ctx.Done():
			cleanup(&wg, wp)
			return nil
		default:
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Fprint(out, ">")
		if err := out.Flush(); err != nil {
			return err
		}

		input, _ := readInput(in)

		switch input {
		case "exit":
			cancel()
			cleanup(&wg, wp)
			return nil
		case "acceptOrder":
			errChan <- controller.WAcceptOrder(ctx, oS, wp, errChan)
		case "returnOrder":
			errChan <- controller.WReturnOrder(ctx, oS, wp, errChan)
		case "placeOrder":
			errChan <- controller.WPlaceOrder(ctx, oS, wp, errChan)
		case "listOrders":
			errChan <- controller.WListOrders(ctx, oS, wp, errChan)
		case "refundOrder":
			errChan <- controller.WRefundOrder(ctx, oS, wp, errChan)
		case "listReturns":
			errChan <- controller.WListReturns(ctx, oS, wp, errChan)
		case "help":
			showHelp()
		case "workers-num":
			errChan <- controller.WChangeNumOfWorkers(ctx, wp, errChan)
		default:
			fmt.Fprintln(out, "Unknown command")
		}
	}
}

func readInput(in *bufio.Reader) (string, error) {
	var input string
	if _, err := fmt.Fscanln(in, &input); err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func listenChannels(ctx context.Context, errChan chan error, notification chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				fmt.Println("Error channel closed")
				return
			}
			if err != nil {
				fmt.Println("\n", err)
			}
		case note, ok := <-notification:
			if !ok {
				fmt.Println("Notification channel closed")
				return
			}
			fmt.Println("\n", note)
		case <-ctx.Done():
			close(errChan)
			return
		}
	}
}

func showHelp() {
	fmt.Println(helpText)
}

func signalSearch(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	for {
		select {
		case <-signals:
			fmt.Printf("\nPress Enter for Exit\n")
			cancel()
		case <-ctx.Done():
			return
		}

	}

}

func cleanup(wg *sync.WaitGroup, wp *wpool.WorkerPool) {
	wp.Stop()
	wg.Wait()
}
