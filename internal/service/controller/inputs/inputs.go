package inputs

import (
	"bufio"
	"errors"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"os"
	"strconv"
	"strings"
	"time"
)

func CollectOrderInput() (*models.Order, string, bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var order models.Order
	var packageType string
	var dateString string

	fmt.Println("Input OrderID _ UserID _ Date(form[2024-12(m)-12(d)])")
	fmt.Print(">")
	if !scanner.Scan() {
		return nil, "", false, fmt.Errorf("Input error: failed to read line")
	}
	input := scanner.Text()
	_, err := fmt.Sscanf(input, "%d %d %s", &order.ID, &order.UserID, &dateString)
	if err != nil {
		return nil, "", false, fmt.Errorf("Input parse Err: %w", err)
	}

	order.KeepUntilDate, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, "", false, fmt.Errorf("Date parse Err: %w", err)
	}

	fmt.Println("Input weight[kg], price[₽], package type [box, bundle, wrap]")
	fmt.Print(">")
	if !scanner.Scan() {
		return nil, "", false, fmt.Errorf("Input error: failed to read line")
	}
	input = scanner.Text()
	_, err = fmt.Sscanf(input, "%d %d %s", &order.Weight, &order.Price, &packageType)
	if err != nil {
		return nil, "", false, fmt.Errorf("Input parse Err: %w", err)
	}

	var answer string
	needWrapping := false
	fmt.Println("Would you like to add a wrap to your package? ['y' - yes, 'n' - no]")
	_, err = fmt.Scan(&answer)
	if err != nil {
		return nil, "", false, fmt.Errorf("Error scaning answer %w", err)
	}
	switch answer {
	case "y":
		needWrapping = true
	case "n":
		needWrapping = false
	default:
		return &models.Order{}, "", false, errors.New("invalid input")
	}

	return &order, packageType, needWrapping, nil

}

func InputOrderID() (uint, error) {
	var id uint
	fmt.Print("Input OrderID\n>")
	_, err := fmt.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error scanning OrderID: %w", err)
	}
	return id, nil
}

func InputUserID() (uint, error) {
	var id uint
	fmt.Print("Input UserID\n>")
	_, err := fmt.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error scanning UserID: %w", err)
	}
	return id, nil
}

func InputOrderIDs() ([]uint, error) {
	fmt.Print("Input all IDs that you want to place\n>")
	reader := bufio.NewReader(os.Stdin)
	temp, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error reading input string: %w", err)
	}

	data := strings.Fields(temp)
	var uintdata []uint
	for _, v := range data {
		uval, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Error converting ID to integer: %w", err)
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
	_, err := fmt.Scan(&choice)
	if err != nil {
		return 0, fmt.Errorf("Error scanning list choice: %w", err)
	}
	return choice, nil
}

func InputN() (int, error) {
	var n int
	fmt.Print("Input n\n>")
	_, err := fmt.Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("Error scanning n: %w", err)
	}
	return n, nil
}

func InputOrderAndUserID() (uint, uint, error) {
	var orderId, userId uint
	fmt.Print("Input OrderID and UserID\n>")
	_, err := fmt.Scan(&orderId, &userId)
	if err != nil {
		return 0, 0, fmt.Errorf("Error scanning OrderID and UserID: %w", err)
	}
	return orderId, userId, nil
}

func InputReturnsPagination() (int, int, error) {
	var limit, page int
	fmt.Print("Input max Returns on one page and Page[1,2,...,n]\n>")
	_, err := fmt.Scan(&limit, &page)
	if err != nil {
		return 0, 0, fmt.Errorf("Error scanning pagination input: %w", err)
	}
	return limit, page, nil
}
