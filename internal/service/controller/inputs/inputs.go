package inputs

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"time"
)

func CollectOrderInput() (*models.Order, string, error) {
	var order models.Order
	var packageType string
	var dateString string

	fmt.Println("Input OrderID _ UserID _ Date(form[2024-12(m)-12(d)])")
	fmt.Print(">")
	_, err := fmt.Scan(&order.ID, &order.UserID, &dateString)
	if err != nil {
		return nil, "", fmt.Errorf("Input api Err: %w\n", err)
	}

	order.KeepUntilDate, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, "", fmt.Errorf("Date parse Err: %w\n", err)
	}

	fmt.Println("Input weight[kg], price[₽], package type [box, bundle, wrap]")
	fmt.Print(">")
	_, err = fmt.Scan(&order.Weight, &order.Price, &packageType)
	if err != nil {
		return nil, "", fmt.Errorf("Input api Err: %w\n", err)
	}

	return &order, packageType, nil
}
