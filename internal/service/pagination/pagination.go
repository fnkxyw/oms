package pagination

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"os"
	"strings"
)

// пагинация скроллом
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
			fmt.Printf("OrderID: %v, Reciver: %v, State: %v, Price: %v₽, Date until which it will be stored: %v ",
				orders[i].ID, orders[i].UserID, orders[i].State, orders[i].Price, orders[i].KeepUntilDate)
		}
		lastIndex = end

		if lastIndex >= total {
			fmt.Println("")
			fmt.Println("End.")
			return nil
		}

		input, _ := reader.ReadString('\n')
		_ = strings.TrimSpace(input)

	}

}

func PagePagination(returns []*models.Order, page, limit int) error {
	if page < 1 || limit < 1 {
		return e.ErrLimitPage
	}

	offset := (page - 1) * limit

	if offset >= len(returns) {
		return e.ErrNoMoreItems
	}

	end := offset + limit
	if end > len(returns) {
		end = len(returns)
	}
	returns = returns[offset:end]
	for _, v := range returns {
		fmt.Printf("OrderID: %v, UserID: %v \n", v.ID, v.UserID)
	}
	return nil
}
