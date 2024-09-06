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
     acceptOrder - allows you to take the order from the user
     returnOrder - allows you to return the order to the user
     placeOrder - allow the order to be released to the user
     listOrders - allows you to get a list of orders  
     returnUser - allows you to accept a return from a user
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
			storage.WriteToJSON("api/orders.json", oS)
			storage.WriteToJSON("api/returns.json", rS)
			return nil
		case "acceptOrder":
			err = service.WAcceptOrder(oS)
			if err != nil {
				fmt.Print(err)
			}
			break
		case "returnOrder":
			err = service.WReturnOrder(oS)
			if err != nil {
				fmt.Print(err)
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
			err = service.WReturnUser(rS, oS)
			if err != nil {
				fmt.Print(err)
			}
			break
		case "listReturns":
			err = service.WListReturns(rS)
			if err != nil {
				fmt.Println(err)
			}
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
