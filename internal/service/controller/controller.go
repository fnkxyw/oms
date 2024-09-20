package controller

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller/inputs"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
	"os"
	"strconv"
	"strings"
)

func WAcceptOrder(s storage.OrderStorageInterface) error {
	order, packageType, err := inputs.CollectOrderInput()
	if err != nil {
		return err
	}

	err = packing.Packing(order, packageType)
	if err != nil {
		return err
	}

	err = orders.AcceptOrder(s, order)
	if err != nil {
		return err
	}

	return nil
}

func WReturnOrder(s storage.OrderStorageInterface) error {
	var id uint
	fmt.Println("Input OrderID")
	fmt.Print(">")
	fmt.Scan(&id)
	if !s.IsConsist(id) {
		return e.ErrNoConsist
	}
	err := orders.ReturnOrder(s, id)
	if err != nil {
		return err
	}
	return nil
}

func WPlaceOrder(s storage.OrderStorageInterface) error {
	fmt.Println("Input all IDs that you want to place")
	fmt.Print(">")
	var temp string
	reader := bufio.NewReader(os.Stdin)
	temp, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	data := strings.Fields(temp)
	var uintdata []uint
	for _, v := range data {

		uval, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Wrong conv id in PlaceOrder")
			return err
		}
		if !s.IsConsist(uint(uval)) {
			return e.ErrNoConsist
		}
		uintdata = append(uintdata, uint(uval))
	}
	err = orders.PlaceOrder(s, uintdata)
	if err != nil {
		return err
	}
	return nil
}

func WListOrders(s storage.OrderStorageInterface) error {
	var (
		id   uint
		n    int
		temp int
	)
	fmt.Println("Input UserID")
	fmt.Print(">")
	fmt.Scan(&id)
	fmt.Println("1.List all orders witch consists on our PuP\n" +
		"2.List last N orders")
	fmt.Print(">")
	fmt.Scan(&temp)
	switch temp {
	case 1:
		err := orders.ListOrders(s, id, 0, true)
		return err
	case 2:
		fmt.Println("Input n")
		fmt.Print(">")
		fmt.Scan(&n)
		err := orders.ListOrders(s, id, n, false)
		return err

	}
	return nil
}

func WRefundOrder(rS storage.ReturnStorageInterface, oS storage.OrderStorageInterface) error {
	fmt.Println("Input OrderID and UserId")
	fmt.Print(">")
	var (
		orderId uint
		userdId uint
	)
	fmt.Scan(&orderId, &userdId)
	err := returns.RefundOrder(rS, oS, orderId, userdId)
	return err
}

func WListReturns(rs storage.ReturnStorageInterface) error {
	fmt.Println("Input max Returns on one page and Page")
	var (
		limit int
		page  int
	)
	fmt.Print(">")
	fmt.Scan(&limit, &page)

	return returns.ListReturns(rs, limit, page)
}
