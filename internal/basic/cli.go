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
     acceptOrder - allows you to take the order from the courier ( form is [acceptOrder OrderID ClientID StorageTime(in hours)] )
     returnOrder - allows you to return the order to the courier ( form is [returnOrder OrderID])
     placeOrder - allow the order to be released to the customer ( form is [placeOrder OrderID1 OrderID2 ... OrderIDn])
     listOrders - allows you to get a list of orders (form is [listOrders ClientID]) 
     returnCustomer - allows you to accept a return from a customer(form is [returnCustomer ClientID OrderID])
     listReturns - allows you to get a list of returns (form is [listReturns ClientID])
`

func Run() error {
	ShowHelp()

	var in *bufio.Reader
	var out *bufio.Writer
	out = bufio.NewWriter(os.Stdout)
	in = bufio.NewReader(os.Stdin)
	fmt.Fprint(out, ">")
	out.Flush()
	strg := storage.NewOrderStorage()

	var input string
	fmt.Fscanln(in, &input)
	for {
		fmt.Fprint(out, ">")
		switch input {
		case "exit":
			storage.WriteToJSON("data/orders.json", strg)
			return nil
		case "acceptOrder":
			service.WAcceptOrder(strg)
			break
		case "returnOrder":
			fmt.Fprint(out, "returnOrder command\n ")
			break
		case "placeOrder":
			fmt.Fprint(out, "placeOrder command\n ")
			break
		case "listOrders":
			fmt.Fprint(out, "listOrders command\n ")
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

		out.Flush()

		fmt.Fscanln(in, &input)
		input = strings.TrimSpace(input)

	}
	return nil
}

func ShowHelp() error {
	fmt.Println(helpText)
	return nil
}
