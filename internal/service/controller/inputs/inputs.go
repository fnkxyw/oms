package inputs

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"os"
	"strconv"
	"strings"
	"time"
)

func CollectOrderInput() (*models.Order, string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var order models.Order
	var packageType string
	var dateString string

	// Считываем первую строку
	fmt.Println("Input OrderID _ UserID _ Date(form[2024-12(m)-12(d)])")
	fmt.Print(">")
	if !scanner.Scan() {
		return nil, "", fmt.Errorf("Input error: failed to read line")
	}
	// Парсим строку
	input := scanner.Text()
	_, err := fmt.Sscanf(input, "%d %d %s", &order.ID, &order.UserID, &dateString)
	if err != nil {
		return nil, "", fmt.Errorf("Input parse Err: %w", err)
	}

	order.KeepUntilDate, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, "", fmt.Errorf("Date parse Err: %w", err)
	}

	fmt.Println("Input weight[kg], price[₽], package type [box, bundle, wrap]")
	fmt.Print(">")
	if !scanner.Scan() {
		return nil, "", fmt.Errorf("Input error: failed to read line")
	}
	// Парсим строку
	input = scanner.Text()
	_, err = fmt.Sscanf(input, "%d %d %s", &order.Weight, &order.Price, &packageType)
	if err != nil {
		return nil, "", fmt.Errorf("Input parse Err: %w", err)
	}

	return &order, packageType, nil
}

func main() {
	// Пример вызова CollectOrderInput()
	order, packageType, err := CollectOrderInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Printf("Order: %+v\n", order)
	fmt.Printf("Package Type: %s\n", packageType)
}

func InputOrderID() (uint, error) {
	var id uint
	fmt.Print("Input OrderID\n>")
	fmt.Scan(&id)
	return id, nil
}

func InputUserID() (uint, error) {
	var id uint
	fmt.Print("Input UserID\n>")
	fmt.Scan(&id)
	return id, nil
}

func InputOrderIDs() ([]uint, error) {
	fmt.Print("Input all IDs that you want to place\n>")
	reader := bufio.NewReader(os.Stdin)
	temp, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	data := strings.Fields(temp)
	var uintdata []uint
	for _, v := range data {
		uval, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Wrong conv id in PlaceOrder")
			return nil, err
		}
		uintdata = append(uintdata, uint(uval))
	}
	return uintdata, nil
}

func InputListChoice() (int, error) {
	var choice int
	fmt.Println("1. List all orderStorage which consists on our PuP\n" +
		"2. List last N orderStorage")
	fmt.Print(">")
	fmt.Scan(&choice)
	return choice, nil
}

func InputN() (int, error) {
	var n int
	fmt.Print("Input n\n>")
	fmt.Scan(&n)
	return n, nil
}

func InputOrderAndUserID() (uint, uint, error) {
	var orderId, userId uint
	fmt.Print("Input OrderID and UserID\n>")
	fmt.Scan(&orderId, &userId)
	return orderId, userId, nil
}

func InputReturnsPagination() (int, int, error) {
	var limit, page int
	fmt.Print("Input max Returns on one page and Page\n>")
	fmt.Scan(&limit, &page)
	return limit, page, nil
}
