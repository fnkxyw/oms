package cli

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
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

func Run(oS *storage.OrderStorage, rS *storage.ReturnStorage) error {
	showHelp()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		fmt.Fprint(out, ">")
		if err := out.Flush(); err != nil {
			return err
		}

		input, err := readInput(in)
		if err != nil {
			return err
		}

		switch input {
		case "exit":
			oS.WritoToJSON()
			rS.WritoToJSON()
			return nil
		case "acceptOrder":
			handleError(service.WAcceptOrder(oS))
		case "returnOrder":
			handleError(service.WReturnOrder(oS))
		case "placeOrder":
			handleError(service.WPlaceOrder(oS))
		case "listOrders":
			handleError(service.WListOrders(oS))
		case "refundOrder":
			handleError(service.WRefundOrder(rS, oS))
		case "listReturns":
			handleError(service.WListReturns(rS))
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

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
func showHelp() error {
	fmt.Println(helpText)
	return nil
}
