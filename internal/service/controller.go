package service

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func WAcceptOrder(s *storage.OrderStorage) error {
	var order models.Order
	fmt.Println("Input OrderID _ UserID _ Date(form[2024-12(m)-12(d)])")

	var dateString string
	_, err := fmt.Scan(&order.ID, &order.UserID, &dateString)
	if err != nil {
		return fmt.Errorf("Input data error: %w", err)
	}

	order.Date, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return fmt.Errorf("Date parse error: %w", err)
	}
	err = AcceptOrder(s, &order)
	if err != nil {
		return err
	}
	return err
}

func WReturnOrder(s *storage.OrderStorage) error {
	var id uint
	fmt.Println("Input OrderID")
	fmt.Print(">")
	fmt.Scan(&id)
	err := ReturnOrder(s, id)
	if err != nil {
		return err
	}
	fmt.Println("Correct!")
	return nil
}

func WPlaceOrder(s *storage.OrderStorage) error {
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
			fmt.Println("Wrond id in PlaceOrder")
			return err
		}
		uintdata = append(uintdata, uint(uval))
	}
	err = PlaceOrder(s, uintdata)
	if err != nil {
		return err
	}
	fmt.Println("Correct!")
	return nil
}

func WListOrders(s *storage.OrderStorage) error {
	var (
		id   uint
		n    int
		temp int
	)
	fmt.Println("Input ClientID")
	fmt.Scan(&id)
	fmt.Println("1.List all orders witch consists on our PuP\n" +
		"2.List last N orders")
	fmt.Scan(&temp)
	switch temp {
	case 1:
		err := ListOrders(s, id, 0, true)
		return err
	case 2:
		fmt.Println("Input n")
		fmt.Scan(&n)
		err := ListOrders(s, id, n, false)
		return err

	}
	return nil
}
