package cli

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"

	"os"
	"strings"
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
`

func Run(oS *orderStorage.OrderStorage) error {
	showHelp()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		fmt.Fprint(out, ">")
		if err := out.Flush(); err != nil {
			return err
		}

		input, _ := readInput(in)

		switch input {
		case "exit":
			err := oS.WriteToJSON()
			if err != nil {
				fmt.Println(err)
			}

			return nil
		case "acceptOrder":
			handleErr(controller.WAcceptOrder(oS))
		case "returnOrder":
			handleErr(controller.WReturnOrder(oS))
		case "placeOrder":
			handleErr(controller.WPlaceOrder(oS))
		case "listOrders":
			handleErr(controller.WListOrders(oS))
		case "refundOrder":
			handleErr(controller.WRefundOrder(oS))
		case "listReturns":
			handleErr(controller.WListReturns(oS))
		case "help":
			showHelp()
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

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Correct!")
	}
}
func showHelp() error {
	fmt.Println(helpText)
	return nil
}
