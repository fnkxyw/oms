package service

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"os"
	"strings"
)

func ScrollPagination(orders []*models.Order, limit int) error {
	total := len(orders)
	lastIndex := 0

	if limit < 1 {
		limit = 1
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press Enter to load next Order")
	for {
		start := lastIndex
		end := lastIndex + limit
		if end > total {
			end = total
		}

		for i := start; i < end; i++ {
			fmt.Printf("OrderID: %v, Reciver: %v, State: %v, Date until which it will be stored: %v \n",
				orders[i].ID, orders[i].UserID, orders[i].State, orders[i].Date)
		}
		lastIndex = end

		if lastIndex >= total {
			fmt.Println("End.")
			return nil
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

	}

	return nil
}
