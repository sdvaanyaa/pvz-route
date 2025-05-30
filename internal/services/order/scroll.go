package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"sort"
)

func (s *orderService) Scroll(userID, lastID string, limit int) ([]*models.Order, string, error) {
	if userID == "" {
		return nil, "", ErrEmptyUserID
	}

	if limit < 1 {
		return nil, "", ErrInvalidLimitNumber
	}

	orders, err := s.storage.GetOrdersByUser(userID)
	if err != nil {
		return nil, "", err
	}

	activeOrders := filterActiveOrders(orders)

	sort.Slice(activeOrders, func(i, j int) bool {
		return activeOrders[i].ID < activeOrders[j].ID
	})

	start := findStartIndex(activeOrders, lastID)

	end := min(start+limit, len(activeOrders))

	selected := activeOrders[start:end]

	nextLastID := getNextLastID(activeOrders, end)

	return selected, nextLastID, nil
}

func findStartIndex(orders []*models.Order, lastID string) int {
	if lastID == "0" {
		return 0
	}
	for i, order := range orders {
		if order.ID == lastID {
			return i + 1
		}
	}
	return 0
}

func getNextLastID(orders []*models.Order, end int) string {
	if end < len(orders) {
		return orders[end-1].ID
	}
	return ""
}
