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
     returnCustomer - allows you to accept a return from a customer
     listReturns - allows you to get a list of returns 
`

func Run(strg *storage.OrderStorage) error {
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
			storage.WriteToJSON("data/orders.json", strg)
			return nil
		case "acceptOrder":
			err = service.WAcceptOrder(strg)
			if err != nil {
				fmt.Println(err)
			}
			break
		case "returnOrder":
			err = service.WReturnOrder(strg)
			if err != nil {
				fmt.Println(err)
			}
			break
		case "placeOrder":
			err = service.WPlaceOrder(strg)
			if err != nil {
				fmt.Print(err)
			}
			break
		case "listOrders":
			err = service.ListOrders(strg, 12, 3, false)
			if err != nil {
				fmt.Println(err)
			}
			break
		case "returnCustomer":
			fmt.Fprint(out, "returnCustomer command\n ")
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
