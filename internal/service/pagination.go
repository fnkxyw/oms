package service

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"os"
	"strings"
)

// пагинация скроллом
func scrollPagination(orders []*models.Order, limit int) error {
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
			fmt.Printf("OrderID: %v, Reciver: %v, State: %v, Date until which it will be stored: %v ",
				orders[i].ID, orders[i].UserID, orders[i].State, orders[i].Date)
		}
		lastIndex = end

		if lastIndex >= total {
			fmt.Println("")
			fmt.Println("End.")
			return nil
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

	}

	return nil
}

// пагинация постраничная
func pagePagination(returns []*models.Return, page, limit int) error {
	if page < 1 || limit < 1 {
		return fmt.Errorf("page and limit must be greater than 0")
	}

	offset := (page - 1) * limit

	if offset >= len(returns) {
		return fmt.Errorf("no more api")
	}

	end := offset + limit
	if end > len(returns) {
		end = len(returns)
	}
	returns = returns[offset:end]
	for _, v := range returns {
		fmt.Printf("OrderID: %v, UserID: %v, Date of return: %v \n", v.ID, v.UserID, v.DateOfReturn)
	}
	return nil
}
