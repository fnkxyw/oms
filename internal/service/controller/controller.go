package controller

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//Файл с обертками для организации входа данных

func WAcceptOrder(s storage.OrderStorageInterface) error {

	var order models.Order
	var pacakgeType string
	fmt.Println("Input OrderID _ UserID _ Date(form[2024-12(m)-12(d)])")
	fmt.Print(">")

	var dateString string
	_, err := fmt.Scan(&order.ID, &order.UserID, &dateString)
	if err != nil {
		return fmt.Errorf("Input api Err: %w\n", err)
	}
	if s.IsConsist(order.ID) {
		return e.ErrIsConsist
	}

	order.KeepUntilDate, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return fmt.Errorf("Date parse Err: %w\n", err)
	}

	fmt.Println("Input weight[kg], price[₽], package type [box, bundle, wrap]")
	fmt.Print(">")
	fmt.Scan(&order.Weight, &order.Price, &pacakgeType)

	err = packing.Packing(&order, pacakgeType)
	if err != nil {
		return err
	}

	err = orders.AcceptOrder(s, &order)
	if err != nil {
		return err
	}

	return err
}

func WReturnOrder(s storage.OrderStorageInterface) error {
	var id uint
	fmt.Println("Input OrderID")
	fmt.Print(">")
	fmt.Scan(&id)
	if !s.IsConsist(id) {
		return e.ErrNoConsist
	}
	err := returns.ReturnOrder(s, id)
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
	err := orders.RefundOrder(rS, oS, orderId, userdId)
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
