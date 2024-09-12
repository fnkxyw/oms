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
			oS.WritoToJSON()
			rS.WritoToJSON()
			return nil
		case "acceptOrder":
			err = service.WAcceptOrder(oS)
			if err != nil {
				fmt.Println(err)
			}

		case "returnOrder":
			err = service.WReturnOrder(oS)
			if err != nil {
				fmt.Println(err)
			}

		case "placeOrder":
			err = service.WPlaceOrder(oS)
			if err != nil {
				fmt.Println(err)
			}

		case "listOrders":
			err = service.WListOrders(oS)
			if err != nil {
				fmt.Println(err)
			}

		case "refundOrder":
			err = service.WRefundOrder(rS, oS)
			if err != nil {
				fmt.Println(err)
			}

		case "listReturns":
			err = service.WListReturns(rS)
			if err != nil {
				fmt.Println(err)
			}

		case "help":
			ShowHelp()

		default:
			fmt.Fprint(out, "Unknown command\n")

		}

		fmt.Fprint(out, ">")
		out.Flush()
		input = ""
		fmt.Fscanln(in, &input)
		input = strings.TrimSpace(input)

	}
}

func ShowHelp() error {
	fmt.Println(helpText)
	return nil
}
