package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"sort"
)

func (s *orderService) Scroll(userID, lastID string, limit int) ([]*models.Order, string, error) {
	const op = "services.order.Scroll"

	if userID == "" {
		return nil, "", fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	if limit < 1 {
		return nil, "", fmt.Errorf("%s: %w", op, ErrInvalidLimitNumber)
	}

	orders, err := s.storage.GetOrdersByUser(userID)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	activeOrders := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status != models.StatusArchived {
			activeOrders = append(activeOrders, order)
		}
	}

	sort.Slice(activeOrders, func(i, j int) bool {
		return activeOrders[i].ID < activeOrders[j].ID
	})

	start := 0
	if lastID != "0" {
		for i, order := range activeOrders {
			if order.ID == lastID {
				start = i + 1
				break
			}
		}
	}

	end := start + limit
	if end > len(activeOrders) {
		end = len(activeOrders)
	}

	selected := activeOrders[start:end]

	nextLastID := ""
	if end < len(activeOrders) {
		nextLastID = activeOrders[end-1].ID
	}

	return selected, nextLastID, nil
}
