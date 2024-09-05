package basic

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
     placeOrder - allow the order to be released to the customer
     listOrders - allows you to get a list of orders  
     returnUser - allows you to accept a return from a customer
     listReturns - allows you to get a list of returns 
`

func Run(oS *storage.OrderStorage, rS *storage.ReturnStorage) error {
	ShowHelp()

	var in *bufio.Reader
	var out *bufio.Writer
	out = bufio.NewWriter(os.Stdout)
	in = bufio.NewReader(os.Stdin)
	fmt.Fprint(out, ">")
	out.Flush()
	var err error
	var input string
	fmt.Fscanln(in, &input)
	for {
		switch input {
		case "exit":
			storage.WriteToJSON("data/orders.json", oS)
			storage.WriteToJSON("data/returns.json", rS)
			return nil
		case "acceptOrder":
			err = service.WAcceptOrder(oS)
			if err != nil {
				fmt.Println(err)
			}
			break
		case "returnOrder":
			err = service.WReturnOrder(oS)
			if err != nil {
				fmt.Println(err)
			}
			break
		case "placeOrder":
			err = service.WPlaceOrder(oS)
			if err != nil {
				fmt.Print(err)
			}
			break
		case "listOrders":
			err = service.WListOrders(oS)
			if err != nil {
				fmt.Println(err)
			}
			break
		case "returnUser":

			break
		case "listReturns":
			fmt.Fprintln(out, "listReturns command\n ")
			break
		case "help":
			ShowHelp()
			break
		default:
			fmt.Fprint(out, "Unknown command\n")
			break
		}

		fmt.Fprint(out, ">")
		out.Flush()
		input = ""
		fmt.Fscanln(in, &input)
		input = strings.TrimSpace(input)

	}
	return nil
}

func ShowHelp() error {
	fmt.Println(helpText)
	return nil
}
